package Test

import (
	"SomeProject/internal/db"
	"SomeProject/internal/handlers"
	"SomeProject/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBanner(t *testing.T) {
	dbPool, err := db.ConnectDB(db.NewConfig())
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer dbPool.Pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/banner", handlers.NewBanner(dbPool)).Methods("POST")

	// Имитируем запрос на создание баннера
	newBanner := models.BannerBase{
		TagIDs:    []int{1, 2},
		FeatureID: 1,
		Content: struct {
			Title string "json:\"title\""
			Text  string "json:\"text\""
			URL   string "json:\"url\""
		}{
			Title: "New Banner Title",
			Text:  "New Banner Text",
			URL:   "http://example.com/new",
		},
		IsActive: true,
	}
	body, _ := json.Marshal(newBanner)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидаемый статус код: %d, получен: %d", http.StatusCreated, w.Code)
	}

}
