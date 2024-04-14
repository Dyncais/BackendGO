// internal/handlers/user_banner.go
package handlers

import (
	"SomeProject/internal/cache"
	"SomeProject/internal/db"
	"SomeProject/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetUserBanner(bannerCache *cache.BannerCache, dbPool *db.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tagID := r.URL.Query().Get("tag_id")
		featureID := r.URL.Query().Get("feature_id")
		useLastRevision := r.URL.Query().Get("use_last_revision") == "false"
		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, models.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		if tagID == "" || featureID == "" {
			http.Error(w, models.ErrWrongData, http.StatusBadRequest)
			return
		}

		cacheKey := tagID + "_" + featureID + "_" + strconv.FormatBool(useLastRevision)

		if banner, found := bannerCache.Get(cacheKey); found {
			respondWithJSON(w, banner)
			return
		}

		banner, err := db.LoadBannerByParams(dbPool.Pool, tagID, featureID, useLastRevision)

		if err != nil {
			if err.Error() == "BannerOff" && token != "admin_token" {
				http.Error(w, models.ErrNoAccess, http.StatusForbidden)
				return
			} else {
				respondWithJSON(w, banner.Content)
				return
			}
			http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		//bannerCache.Set(cacheKey, banner)

		respondWithJSON(w, banner.Content)
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
