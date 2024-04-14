package getbanner

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"

	"AvitoTech/internal/models"
)

type Handler struct {
	userChecker UserChecker
	DB          DBInt
}

func NewHandler(db DBInt, checker UserChecker) *Handler {
	return &Handler{
		DB:          db,
		userChecker: checker,
	}
}

func (h *Handler) GetBanners(c *fiber.Ctx) error {
	var data []models.Banner

	tagID := c.QueryInt("tag_id", -1)
	featureID := c.QueryInt("feature_id", -1)
	limit := c.QueryInt("limit", -1)
	offset := c.QueryInt("offset", 0)

	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}
	if !admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if tagID == -1 && featureID != -1 {
		data, err = h.DB.GetBannersByFeatureID(uint64(featureID), limit, offset)
	}
	if tagID != -1 && featureID == -1 {
		data, err = h.DB.GetBannersByTagID(uint64(tagID), limit, offset)
	}
	if tagID != -1 && featureID != -1 {
		data, err = h.DB.GetBanners(uint64(featureID), uint64(tagID), limit, offset)
	}
	if tagID == -1 && featureID == -1 {
		return c.JSON(fiber.Map{
			"error": "Missing tag or feature",
		})
	}

	if err != nil {
		log.Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(data)
}
