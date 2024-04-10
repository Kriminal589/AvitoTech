package models

import "github.com/jackc/pgtype"

type DeleteBanner struct {
	ID uint64 `params:"id" json:"id"`
}

type UserBanner struct {
	TagID     string `params:"tag_id" json:"tag_id"`
	FeatureID string `params:"feature_id" json:"feature_id"`
}

type RequestBanner struct {
	FeatureID uint64       `params:"feature_id" json:"feature_id"`
	TagIDS    []byte       `params:"tag_ids" json:"tag_ids"`
	Content   pgtype.JSONB `params:"content" json:"content"`
	IsActive  bool         `params:"is_active" json:"is_active"`
}
