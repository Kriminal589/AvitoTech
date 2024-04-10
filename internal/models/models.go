package models

import (
	"github.com/jackc/pgtype"
	"time"
)

type Tag struct {
	ID uint64 `json:"id"`
}

type Feature struct {
	ID uint64 `json:"id"`
}

type User struct {
	ID        uint64  `json:"id"`
	Admin     bool    `json:"admin"`
	FeatureID uint64  `json:"-"`
	Feature   Feature `json:"Feature"`
	TagID     uint64  `json:"-"`
	Tag       Tag     `json:"Tag"`
}

type Banner struct {
	BannerID  uint64       `json:"banner_id"`
	TagIDS    []uint64     `json:"tag_ids"`
	FeatureID uint64       `json:"feature_id"`
	Content   pgtype.JSONB `json:"content"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
