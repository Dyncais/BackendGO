package main

import (
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"SomeProject/internal/models"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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
	}

	newBanner := models.Banner{
		Title:     "Новый баннер",
		Text:      "Описание нового баннера",
		URL:       "https://example.com",
		TagIDs:    []int{1, 2, 3},
		FeatureID: 1,
		IsActive:  true,
	}

	bannerID, err := db.InsertBanner(dbPool.Pool, newBanner)
	if err != nil {
		log.Fatalf("Ошибка при вставке баннера: %v", err)
	}

	log.Printf("Баннер успешно добавлен с ID %d", bannerID)

	bannerCache := cache.NewBannerCache()

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserBanner(bannerCache, dbPool, w, r)
	}).Methods("GET")

	//изменить логику
	http.HandleFunc("/banners", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodGet {
			//handlers.GetBanner(w, r)
		} else {
			//handlers.GetBanner(w, r)
		}
	})

	router.HandleFunc("/banner/{id}", handlers.PatchBanner(dbPool)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", handlers.DeleteBanner(dbPool)).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
