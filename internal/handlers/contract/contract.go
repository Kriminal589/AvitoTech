package contract

import (
	"AvitoTech/internal/models"
	"github.com/jackc/pgtype"
)

type DBInt interface {
	GetUserBanner(tagId uint64, featureId uint64) (pgtype.JSONB, error)
	GetBannersWithFeatureId(featureId uint64) ([]models.Banner, error)
	DeleteBanner(id uint64) (string, error)
	PostBanner(banner models.RequestBanner) (uint64, error)
	PatchBanner(banner models.RequestBanner, id int) (string, error)
}
