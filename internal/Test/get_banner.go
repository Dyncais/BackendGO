package Test

import (
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBanners(t *testing.T) {
	cfg := db.NewConfig()
	dbPool, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/banner", handlers.GetBanner(dbPool)).Methods("GET")

	// Тестирование без токена
	request, _ := http.NewRequest("GET", "/banner", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Тестирование с токеном
	request.Header.Set("token", "valid_token")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusForbidden {
		t.Errorf("Handler returned wrong status code with valid token: got %v want %v", status, http.StatusOK)
	}

	// Тестирование с токеном админа
	request.Header.Set("token", "admin_token")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code with valid token: got %v want %v", status, http.StatusOK)
	}

}
