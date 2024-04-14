package models

import (
	"time"

	"github.com/jackc/pgtype"
)

// TODO: заменить JSONB на []byte
type Banner struct {
	IsActive  bool         `params:"is_active" json:"is_active"`
	BannerID  uint64       `params:"banner_id" json:"banner_id"`
	TagIDs    []uint64     `params:"tag_ids" json:"tag_ids"`
	FeatureID uint64       `params:"feature_id" json:"feature_id"`
	Content   pgtype.JSONB `params:"content" json:"content"`
	CreatedAt time.Time    `params:"created_at" json:"created_at"`
	UpdatedAt time.Time    `params:"updated_at" json:"updated_at"`
}
