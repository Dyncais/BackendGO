package Test

import (
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteBanner(t *testing.T) {
	cfg := db.NewConfig()
	dbPool, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/banner/{id}", handlers.DeleteBanner(dbPool)).Methods("DELETE")

	// Попытка удаления баннера без токена
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/banner/1"), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Попытка удаления баннера с неверным токеном
	request.Header.Set("token", "invalid_token")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusForbidden {
		log.Println("2")
		t.Errorf("Handler returned wrong status code with invalid token: got %v want %v", status, http.StatusForbidden)
	}

	// Попытка удаления баннера с верным токеном
	request.Header.Set("token", "admin_token")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code with valid token: got %v want %v", status, http.StatusNoContent)
	}

	if err := DropTables(dbPool.Pool); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

}
