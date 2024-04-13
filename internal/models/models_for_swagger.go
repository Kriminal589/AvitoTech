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

type TokenMetadata struct {
	Id uint64 `params:"id" json:"id"`
}

type BannerTagLink struct {
	BannerID uint64 `params:"banner_id" json:"banner_id"`
	TagID    uint64 `params:"tag_id" json:"tag_id"`
}
