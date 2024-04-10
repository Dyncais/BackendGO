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
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid banner ID", http.StatusBadRequest)
			return
		}

		var updateData models.BannerUpdate
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := db.UpdateBanner(dbPool.Pool, id, updateData); err != nil {
			http.Error(w, "Failed to update banner", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
func DeleteBanner(dbPool *db.DBPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("DeleteBanner called")
		vars := mux.Vars(r)
		idStr := vars["id"]
		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, "Unauthorized: Token is required", http.StatusUnauthorized)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid banner ID", http.StatusBadRequest)
			panic(err)
		}

		err = db.DeleteBanner(dbPool.Pool, id)
		if err != nil {
			log.Printf("Failed to delete banner with ID %d: %v", id, err)
			if err.Error() == "no banner found with ID "+strconv.Itoa(id) {
				http.Error(w, "Banner not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
