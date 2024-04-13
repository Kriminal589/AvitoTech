package handlers

import (
	"AvitoTech/internal/models"
	"github.com/gofiber/fiber/v2"
)

type DBInt interface {
	GetBannersByFeatureId(featureId uint64, limit int, offset int) ([]models.Banner, error)
	GetBannersByTagId(tagId uint64, limit int, offset int) ([]models.Banner, error)
	GetBanners(featureId uint64, tagId uint64, limit int, offset int) ([]models.Banner, error)
	DeleteBanner(id uint64) error
	PostBanner(banner models.RequestBanner) (uint64, error)
	PatchBanner(banner models.RequestBanner, id int) error
	GetUserRole(id uint64) (bool, error)
}

type BannerGetter interface {
	GetUserBanner(tagId uint64, featureId uint64, lastRevision bool, isAdmin bool) ([]byte, error)
}

type UserChecker interface {
	IsAdmin(c *fiber.Ctx) (bool, error)
}
