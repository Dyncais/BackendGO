// /internal/models/banner.go
package models

type Banner struct {
	ID        int32 `json:id`
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_Active"`
}

type BannerWithID struct {
	ID        int32  `json:id`
	Title     string `json:"title"`
	Text      string `json:"text"`
	URL       string `json:"url"`
	IsActive  bool   `json:"is_Active"`
	CreatedAt string `json:"created_At"`
	UpdatedAt string `json:"updated_At"`
}

type BannerUpdate struct {
	Title     *string `json:"title,omitempty"`
	Text      *string `json:"text,omitempty"`
	URL       *string `json:"url,omitempty"`
	TagIDs    *[]int  `json:"tag_ids,omitempty"`
	FeatureID *int    `json:"feature_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type BannerCreationRequest struct {
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_active"`
}

/*type BannerResponse struct {
	Content Content `json:"content"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}*/
