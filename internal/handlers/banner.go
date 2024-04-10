package handlers

/*func GetFilteredBanner(bannerCache *cache.BannerCache, dbPool *db.DBPool, w http.ResponseWriter, r *http.Request) {
	log.Println("GetFilteredBanner called")
	token := r.Header.Get("token")

	if token == "" {
		http.Error(w, "Unauthorized: Token is required", http.StatusUnauthorized)
		return
	}

	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	banners, err := db.GetBanners(dbPool.Pool, featureID, tagID, limit, offset)
	if err != nil {
		log.Printf("Ошибка при получении баннеров: %v\n", err)
		http.Error(w, "Failed to get banners", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

type BannerRequest struct {
	TagIDs    []int             `json:"tag_ids"`
	FeatureID int               `json:"feature_id"`
	Content   map[string]string `json:"content"`
	IsActive  bool              `json:"is_active"`
}

func PostFilteredBanner(dbPool *db.DBPool, w http.ResponseWriter, r *http.Request) {
	log.Println("PostFilteredBanner called")
	token := r.Header.Get("token")

	if token == "" {
		http.Error(w, "Unauthorized: Token is required", http.StatusUnauthorized)
		return
	}

	var req BannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
		return
	}

	bannerID, err := db.CreateBanner(dbPool.Pool, req)
	if err != nil {
		log.Printf("Ошибка при создании баннера: %v\n", err)
		http.Error(w, "Failed to create banner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"banner_id": bannerID})
}
*/
