package models

import (
	"time"
)

type ModelMap map[string]interface{}

type BannerExpanded struct {
	ID        int32     `json:"banner_id,omitempty"`
	TagIds    []int32   `json:"tag_ids,omitempty"`
	FeatureId int32     `json:"feature_id,omitempty"`
	Content   ModelMap  `json:"content,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type IdResponse struct {
	// Идентификатор созданного баннера
	BannerId int32 `json:"banner_id,omitempty"`
}

type Banner struct {
	ID        int32    `json:"id"`
	TagIds    []int32  `json:"tag_ids"`
	FeatureId int32    `json:"feature_id"`
	Content   ModelMap `json:"content"`
	IsActive  bool     `json:"is_active"`
}

type BannerNoId struct {
	TagIds    []int32  `json:"tag_ids,omitempty"`
	FeatureId int32    `json:"feature_id,omitempty"`
	Content   ModelMap `json:"content,omitempty"`
	IsActive  bool     `json:"is_active,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
