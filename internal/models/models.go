package models

type Tag struct {
	ID uint64 `json:"id"`
}

type Feature struct {
	ID uint64 `json:"id"`
}

type User struct {
	ID        uint64  `json:"id"`
	Admin     bool    `json:"Role"`
	FeatureID uint64  `json:"-"`
	Feature   Feature `json:"Feature"`
	TagID     uint64  `json:"-"`
	Tag       Tag     `json:"Tag"`
}

type Banner struct {
	ID        uint64  `json:"id"`
	FeatureID uint64  `json:"-"`
	Feature   Feature `json:"Feature"`
	TagID     uint64  `json:"-"`
	Tag       Tag     `json:"Tag"`
}
