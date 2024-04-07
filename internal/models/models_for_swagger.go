package models

type DeleteBanner struct {
	ID uint64 `params:"id" json:"id"`
}

type UserBanner struct {
	TagID     string `params:"tag_id" json:"tag_id"`
	FeatureID string `params:"feature_id" json:"feature_id"`
}
