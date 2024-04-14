package getbanner

import (
	"github.com/gofiber/fiber/v2"

	"AvitoTech/internal/models"
)

//go:generate mockery --name DBInt --with-expecter
type DBInt interface {
	GetBannersByFeatureID(featureID uint64, limit int, offset int) ([]models.Banner, error)
	GetBannersByTagID(tagID uint64, limit int, offset int) ([]models.Banner, error)
	GetBanners(featureID uint64, tagID uint64, limit int, offset int) ([]models.Banner, error)
}

//go:generate mockery --name UserChecker --with-expecter
type UserChecker interface {
	IsAdmin(c *fiber.Ctx) (bool, error)
}
