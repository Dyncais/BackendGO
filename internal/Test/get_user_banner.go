package Test

import (
	"SomeProject/internal/cache"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"github.com/gorilla/mux"
)

func TestGetUserBanner(t *testing.T) {
	cfg := db.NewConfig()
	dbPool, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Pool.Close()

	if err := CreateTables(dbPool.Pool); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	if err := InsertTestData(dbPool.Pool); err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	router := mux.NewRouter()
	bannerCache := cache.NewBannerCache()
	router.HandleFunc("/user_banner", handlers.GetUserBanner(bannerCache, dbPool)).Methods("GET")

	// Тестирование без токена
	request, _ := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Тестирование с токеном
	request.Header.Set("token", "valid_token")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code with valid token: got %v want %v", status, http.StatusOK)
	}

}
