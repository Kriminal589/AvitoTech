package handlers

import (
	"github.com/gofiber/fiber/v2"

	"AvitoTech/internal/models"
)

type DBInt interface {
	GetBannersByFeatureID(featureID uint64, limit int, offset int) ([]models.Banner, error)
	GetBannersByTagID(tagID uint64, limit int, offset int) ([]models.Banner, error)
	GetBanners(featureID uint64, tagID uint64, limit int, offset int) ([]models.Banner, error)
	DeleteBanner(id uint64) error
	PostBanner(banner models.Banner) (uint64, error)
	PatchBanner(banner models.Banner, id int) error
	GetUserRole(id uint64) (bool, error)
}

type BannerGetter interface {
	GetUserBanner(tagID uint64, featureID uint64, lastRevision bool, isAdmin bool) ([]byte, error)
}

type UserChecker interface {
	IsAdmin(c *fiber.Ctx) (bool, error)
}
