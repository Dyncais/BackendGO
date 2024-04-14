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

func strPtr(s string) *string {
	return &s
}
func TestPatchBanner(t *testing.T) {
	dbPool, err := db.ConnectDB(db.NewConfig())
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer dbPool.Pool.Close()

	// Очистка и подготовка тестовой БД
	//test.Setup(dbPool) // Предполагается, что эта функция создает таблицы и заполняет тестовыми данными

	// Создаем http сервер с нашими обработчиками
	router := mux.NewRouter()
	router.HandleFunc("/banner/{id}", handlers.PatchBanner(dbPool)).Methods("PATCH")

	// Имитируем запрос на обновление баннера
	updateData := models.BannerUpdate{
		Title: strPtr("Updated Title"),
		Text:  strPtr("Updated Text"),
		URL:   strPtr("http://example.com/updated"),
	}
	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PATCH", "/banner/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, получен: %d", http.StatusOK, w.Code)
	}

	// Проверка, что баннер действительно обновлен, может быть добавлена здесь
}
