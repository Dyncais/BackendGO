// /internal/models/banner.go
package models

import "time"

type BannerInfo struct {
	Content struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
}

type BannerBase struct {
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive bool `json:"is_active"`
}

type Banner struct {
	BannerBase
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerCreationRequest struct {
	BannerBase
}

type BannerUpdate struct {
	Title     *string `json:"title,omitempty"`
	Text      *string `json:"text,omitempty"`
	URL       *string `json:"url,omitempty"`
	TagIDs    *[]int  `json:"tag_ids,omitempty"`
	FeatureID *int    `json:"feature_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
