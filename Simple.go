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

	if err := Test.DropTables(dbPool.Pool); err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	if err := Test.CreateTables(dbPool.Pool); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	if err := Test.InsertTestData(dbPool.Pool); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	bannerCache := cache.NewBannerCache()

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", handlers.GetUserBanner(bannerCache, dbPool)).Methods("GET")

	router.HandleFunc("/banner", handlers.GetBanner(dbPool)).Methods("GET")
	router.HandleFunc("/banner", handlers.NewBanner(dbPool)).Methods("POST")

	router.HandleFunc("/banner/{id}", handlers.PatchBanner(dbPool)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", handlers.DeleteBanner(dbPool)).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
