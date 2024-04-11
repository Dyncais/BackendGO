package main

import (
	"SomeProject/internal/Test"
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

/*func createTables(dbPool *pgxpool.Pool) error {

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

	return nil
}

func insertTestData(dbPool *pgxpool.Pool) error {

	_, err := dbPool.Exec(context.Background(), `
        INSERT INTO tags (name) VALUES ('Тег 1'), ('Тег 2'), ('Тег 3')
        ON CONFLICT DO NOTHING;`)
	if err != nil {
		return err
	}

	_, err = dbPool.Exec(context.Background(), `
        INSERT INTO features (name) VALUES ('Фича 1'), ('Фича 2'), ('Фича 3')
        ON CONFLICT DO NOTHING;`)
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

	return nil
}
*/

func main() {

	cfg := db.NewConfig()

	if err := db.CreateDatabase(cfg); err != nil {
		log.Printf("Warning: Unable to create database: %v", err)
	}

	dbPool, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Pool.Close()

	/*if err := createTables(dbPool.Pool); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	if err := insertTestData(dbPool.Pool); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	createTableQuery := `
	    CREATE TABLE IF NOT EXISTS banners1 (
	        id SERIAL PRIMARY KEY,
	        Title VARCHAR(255) NOT NULL,
	        Text TEXT,
	        URL VARCHAR(255),
	        TagIDs INTEGER[],
	        FeatureID INT,
	        IsActive BOOLEAN NOT NULL DEFAULT TRUE
	    );`

		_, err = dbPool.Pool.Exec(context.Background(), createTableQuery)
		if err != nil {
			log.Fatalf("Ошибка при создании таблицы banners: %v", err)
		}*/

	if err := Test.DropTables(dbPool.Pool); err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	// Вызов функции для создания таблиц
	if err := Test.CreateTables(dbPool.Pool); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Вызов функции для вставки тестовых данных
	if err := Test.InsertTestData(dbPool.Pool); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	bannerCache := cache.NewBannerCache()

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", handlers.GetUserBanner(bannerCache, dbPool)).Methods("GET")

	//router.HandleFunc("/banner", handlers.GetBanner(dbPool)).Methods("GET")
	router.HandleFunc("/banner", handlers.NewBanner(dbPool)).Methods("POST")

	router.HandleFunc("/banner/{id}", handlers.PatchBanner(dbPool)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", handlers.DeleteBanner(dbPool)).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
