package models

type DeleteBanner struct {
	ID uint64 `params:"id" json:"id"`
}

type UserBanner struct {
	TagID     uint64 `json:"tag_id"`
	FeatureID uint64 `json:"feature_id"`
}
