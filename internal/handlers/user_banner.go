// internal/handlers/user_banner.go
package handlers

import (
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"SomeProject/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetUserBanner(bannerCache *cache.BannerCache, dbPool *db.DBPool, w http.ResponseWriter, r *http.Request) {
	log.Println("GetUserBanner called")
	tagID := r.URL.Query().Get("tag_id")
	featureID := r.URL.Query().Get("feature_id")
	useLastRevision := r.URL.Query().Get("use_last_revision") == "false"
	token := r.Header.Get("token")

	if token == "" {
		http.Error(w, "Unauthorized: Token is required", http.StatusUnauthorized)
		return
	}

	// Проверка наличия обязательных параметров
	if tagID == "" || featureID == "" {
		http.Error(w, "Bad Request: tag_id and feature_id are required", http.StatusBadRequest)
		return
	}

	cacheKey := tagID + "_" + featureID + "_" + strconv.FormatBool(useLastRevision)

	// Поиск в кэше
	if banner, found := bannerCache.Get(cacheKey); found {
		respondWithJSON(w, banner)
		return
	}

	banner, err := db.LoadBannerByParams(dbPool.Pool, tagID, featureID, useLastRevision)
	if err != nil {
		http.Error(w, "Failed to load banner: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bannerCache.Set(cacheKey, banner)

	respondWithJSON(w, banner)
}

func respondWithJSON(w http.ResponseWriter, data models.Banner) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
