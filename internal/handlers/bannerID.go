package handlers

import (
	"SomeProject/internal/db"
	"SomeProject/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type BannerUpdateRequest struct {
	TagIDs    []int             `json:"tag_ids,omitempty"`
	FeatureID int               `json:"feature_id,omitempty"`
	Content   map[string]string `json:"content,omitempty"`
	IsActive  bool              `json:"is_active,omitempty"`
}

func PatchBanner(dbPool *db.DBPool) func(w http.ResponseWriter, r *http.Request) {
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

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, models.ErrFailedLoadBanner, http.StatusNotFound)
			return
		}

		var updateData models.BannerUpdate
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, models.ErrWrongData, http.StatusBadRequest)
			return
		}

		if err := db.UpdateBanner(dbPool.Pool, id, updateData); err != nil {
			http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
func DeleteBanner(dbPool *db.DBPool) func(w http.ResponseWriter, r *http.Request) {
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

		log.Println("DeleteBanner called")
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, models.ErrWrongData, http.StatusNotFound)
			panic(err)
		}

		err = db.DeleteBanner(dbPool.Pool, id)
		if err != nil {
			if err.Error() == "no banner found with ID "+strconv.Itoa(id) {
				http.Error(w, models.ErrFailedLoadBanner, http.StatusNotFound)
			} else {
				http.Error(w, models.ErrInternalServerError, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
