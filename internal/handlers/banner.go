package handlers

import (
	"SomeProject/internal/db"
	"SomeProject/internal/models"
	"encoding/json"
	"net/http"
)

func GetBanner(dbPool *db.DBPool) http.HandlerFunc {
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

		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if limit == "" {
			limit = "10"
		}
		if offset == "" {
			offset = "0"
		}

		banners, err := db.LoadBannersByParams(dbPool.Pool, tagID, featureID, limit, offset)
		if err != nil {
			http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, banners)
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
