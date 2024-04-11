// internal/handlers/user_banner.go
package handlers

import (
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetUserBanner(bannerCache *cache.BannerCache, dbPool *db.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetUserBanner called")
		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		useLastRevision := r.URL.Query().Get("use_last_revision") == "true" // Предполагается, что "true" по умолчанию
		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, "Unauthorized: Token is required", http.StatusUnauthorized)
			return
		}

		if tagID == "" || featureID == "" {
			http.Error(w, "Bad Request: tag_id and feature_id are required", http.StatusBadRequest)
			return
		}

		cacheKey := tagID + "_" + featureID + "_" + strconv.FormatBool(useLastRevision)

		if banner, found := bannerCache.Get(cacheKey); found {
			respondWithJSON(w, banner)
			return
		}

		banner, err := db.LoadBannerByParams(dbPool.Pool, tagID, featureID, useLastRevision)
		if err != nil {
			http.Error(w, "Failed to load banner: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//bannerCache.Set(cacheKey, banner)

		respondWithJSON(w, banner)
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
