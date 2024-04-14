package Test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func DropTables(dbPool *pgxpool.Pool) error {
	log.Printf("Начало дропа")
	tables := []string{"banner_tags", "banners1", "tags", "features", "banners"}
	for _, table := range tables {
		_, err := dbPool.Exec(context.Background(), "DROP TABLE IF EXISTS "+table+" CASCADE;")
		if err != nil {
			return err
		}
	}
	log.Printf("Дропнул")
	return nil
}

func CreateTables(dbPool *pgxpool.Pool) error {
	log.Printf("Начало создание")
	_, err := dbPool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS tags (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS features (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );`)
	if err != nil {
		return err
	}

	createFunction := `
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`

	_, err = dbPool.Exec(context.Background(), createFunction)
	if err != nil {
		log.Fatalf("Ошибка при создании функции для обновления: %v", err)
	}

	createTableQuery := `
CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Text TEXT,
    URL VARCHAR(255),
    TagIDs INTEGER[],
    FeatureID INT,
    IsActive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
`

	_, err = dbPool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы banners: %v", err)
	}

	createTrigger := `
CREATE TRIGGER update_banners_updated_at BEFORE UPDATE ON banners
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
`

	_, err = dbPool.Exec(context.Background(), createTrigger)
	if err != nil {
		log.Fatalf("Ошибка при создании триггера: %v", err)
	}

	_, err = dbPool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS banner_tags (
    	banner_id INTEGER NOT NULL,
    	tag_id INTEGER NOT NULL,
    	CONSTRAINT fk_banner FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    	CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    	PRIMARY KEY (banner_id, tag_id));
	`)
	if err != nil {
		return err
	}

	log.Printf("Создал")
	return nil
}

func InsertTestData(dbPool *pgxpool.Pool) error {
	log.Printf("Начало вставки")
	_, err := dbPool.Exec(context.Background(), `
        INSERT INTO tags (name) VALUES ('Тег 1'), ('Тег 2'), ('Тег 3'), ('Тег 4'), ('Тег 5')
        ON CONFLICT DO NOTHING;`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(context.Background(), `
        INSERT INTO features (name) VALUES ('Фича 1'), ('Фича 2'), ('Фича 3'), ('Фича 4'), ('Фича 5')
        ON CONFLICT DO NOTHING;`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(context.Background(), `
	INSERT INTO banners (Title, Text, URL, TagIDs, FeatureID, IsActive) VALUES 
	('Баннер 1', 'Текст баннера 1', 'http://example.com/1', ARRAY[1, 2], 1, TRUE),
	('Баннер 2', 'Текст баннера 2', 'http://example.com/2', ARRAY[3, 4], 2, TRUE)
	ON CONFLICT (id) DO NOTHING;`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(context.Background(), `
        INSERT INTO banner_tags (banner_id, tag_id) VALUES 
		(1, 1), 
		(1, 2);`)
	if err != nil {
		return err
	}
	log.Printf("Вставил")
	return nil
}

func TestAPIPerformance() (int64, float64, float64, int64) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
		},
	}
	var count atomic.Int64
	var errorCount atomic.Int64
	var wg sync.WaitGroup

	duration := 5 * time.Second
	startTime := time.Now()
	endTime := startTime.Add(duration)

	for time.Now().Before(endTime) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tagID := rand.Intn(1000)
			featureID := rand.Intn(1000)

			url := fmt.Sprintf("http://localhost:8080/user_banner?tag_id=%d&feature_id=%d", tagID, featureID)
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				errorCount.Add(1)
				return
			}

			resp, err := client.Do(request)
			if err != nil {
				fmt.Printf("Error executing request: %v\n", err)
				errorCount.Add(1)
				return
			}
			resp.Body.Close()

			count.Add(1)
		}()
		// Добавляем небольшую задержку, чтобы предотвратить чрезмерное создание горутин
		time.Sleep(750 * time.Nanosecond)
	}

	wg.Wait() // Дожидаемся окончания всех горутин

	// Расчет среднего количества запросов в секунду
	elapsed := time.Since(startTime).Seconds()
	averagePerSecond := float64(count.Load()) / elapsed
	return count.Load(), averagePerSecond, elapsed, errorCount.Load()
}

func TestAPIHandle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		count, averagePerSecond, elapsed, errors := TestAPIPerformance()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Удачных запросов: %d, Ошибок %d/%d (%.2f proc), Прошло времени: %f, Запросов в секунду в среднем: %.2f", count, errors, count+errors, float64(errors)/(float64(count+errors)/100), elapsed, averagePerSecond)))
	}
}
