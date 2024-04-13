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
	BannerID  uint64       `params:"banner_id" json:"banner_id"`
	TagIDS    []uint64     `params:"tag_ids" json:"tag_ids"`
	FeatureID uint64       `params:"feature_id" json:"feature_id"`
	Content   pgtype.JSONB `params:"content" json:"content"`
	IsActive  bool         `params:"is_active" json:"is_active"`
	CreatedAt time.Time    `params:"created_at" json:"created_at"`
	UpdatedAt time.Time    `params:"updated_at" json:"updated_at"`
}
