package main

import (
	"SomeProject/internal/Test"
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"
	"time"
)

func main() {

	cfg := db.NewConfig()

	if err := db.CreateDatabase(cfg); err != nil {
		log.Printf("Не удалось создать БД %v", err)
	}

	var dbPool *db.DBPool
	var err error
	for i := 0; i < 15; i++ {
		dbPool, err = db.ConnectDB(cfg)
		if err == nil {
			fmt.Println("Успешно подключено к базе данных")
			break
		}
		fmt.Printf("Попытка %d: Не удалось подключиться к базе данных, ошибка: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond) // Ожидание 100 мс перед следующей попыткой
	}
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Pool.Close()

	const enableTesting = false

	if enableTesting {
		StartTesting(dbPool)
	}

	bannerCache := cache.NewBannerCache()

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", handlers.GetUserBanner(bannerCache, dbPool)).Methods("GET")

	router.HandleFunc("/banner", handlers.GetBanner(dbPool)).Methods("GET")
	router.HandleFunc("/banner", handlers.NewBanner(dbPool)).Methods("POST")

	router.HandleFunc("/banner/{id}", handlers.PatchBanner(dbPool)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", handlers.DeleteBanner(dbPool)).Methods("DELETE")

	if enableTesting {
		router.HandleFunc("/test", Test.TestAPIHandle()).Methods("GET")
	}
	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func StartTesting(dbPool *db.DBPool) {
	if err := Test.DropTables(dbPool.Pool); err != nil {
		log.Fatalf("Не удалось дропнуть: %v", err)
	}

	testing.Init()
	t := new(testing.T)
	Test.TestGetUserBanner(t)
	if t.Failed() {
		log.Println("Тесты на GET /user_banner не пройдены.")
	} else {
		log.Println("Все тесты GET /user_banner успешно пройдены.")
	}
	Test.TestGetBanners(t)
	if t.Failed() {
		log.Println("Тесты на GET /banner не пройдены.")
	} else {
		log.Println("Все тесты GET /banner успешно пройдены.")
	}
	Test.TestCreateBanner(t)
	if t.Failed() {
		log.Println("Тесты на POST /banner не пройдены.")
	} else {
		log.Println("Все тесты POST /banner успешно пройдены.")
	}

	Test.TestPatchBanner(t)
	if t.Failed() {
		log.Println("Тесты на PATCH /banner/{id} не пройдены.")
	} else {
		log.Println("Все тесты PATCH /banner/{id} успешно пройдены.")
	}

	Test.TestDeleteBanner(t)
	if t.Failed() {
		log.Println("Тесты на DELETE /banner/{id} не пройдены.")
	} else {
		log.Println("Все тесты DELETE /banner/{id} успешно пройдены.")
	}
}
