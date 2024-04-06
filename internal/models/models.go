package models

type Tag struct {
	ID uint64 `json:"id"`
}

type Feature struct {
	ID uint64 `json:"id"`
}

type User struct {
	ID      uint64  `json:"id"`
	Feature Feature `json:"Feature"`
	Tag     Tag     `json:"Tag"`
}

type Banner struct {
	ID      uint64  `json:"id"`
	Feature Feature `json:"Feature"`
	Tag     Tag     `json:"Tag"`
}

type HttpParams struct {
	ID uint64 `params:"id" json:"id"`
}
