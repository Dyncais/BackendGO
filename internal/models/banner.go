// /internal/models/banner.go
package models

type Banner struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	URL       string `json:"url"`
	TagIDs    []int  `json:"tag_ids"`
	FeatureID int    `json:"feature_id"`
	IsActive  bool   `json:"is_active"`
}

type BannerUpdate struct {
	Title     *string `json:"title,omitempty"`
	Text      *string `json:"text,omitempty"`
	URL       *string `json:"url,omitempty"`
	TagIDs    *[]int  `json:"tag_ids,omitempty"`
	FeatureID *int    `json:"feature_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

/*type BannerResponse struct {
	Content Content `json:"content"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}*/
