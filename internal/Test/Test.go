package Test

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
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

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS banners (
        id SERIAL PRIMARY KEY,
        Title VARCHAR(255) NOT NULL,
        Text TEXT,
        URL VARCHAR(255),
        TagIDs INTEGER[],
        FeatureID INT,
        IsActive BOOLEAN NOT NULL DEFAULT TRUE
    );`

	_, err = dbPool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы banners: %v", err)
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
