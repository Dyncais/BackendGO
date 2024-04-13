package handlers

import (
	"SomeProject/internal/db"
	"SomeProject/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetBanner(dbPool *db.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetBanner called")

		tagID, featureID := r.URL.Query().Get("tag_id"), r.URL.Query().Get("feature_id")
		limit, offset := r.URL.Query().Get("limit"), r.URL.Query().Get("offset")
		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, models.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		if token != "admin_token" {
			http.Error(w, models.ErrNoAccess, http.StatusForbidden)
			return
		}

		if limit == "" {
			limit = "10"
		}
		if offset == "" {
			offset = "0"
		}

		banner, err := db.LoadBannersByParams(dbPool.Pool, tagID, featureID, limit, offset)
		if err != nil {
			http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		//bannerCache.Set(cacheKey, banner)

		respondWithJSON(w, banner)
	}
}

func NewBanner(dbPool *db.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, models.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		if token != "admin_token" {
			http.Error(w, models.ErrNoAccess, http.StatusForbidden)
			return
		}

		var bannerData models.BannerCreationRequest
		if err := json.NewDecoder(r.Body).Decode(&bannerData); err != nil {
			http.Error(w, models.ErrWrongData, http.StatusBadRequest)
			return
		}

		bannerID, err := db.InsertBanner(dbPool.Pool, bannerData)
		if err != nil {
			http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"banner_id": bannerID})
	}
}
