package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"

	"AvitoTech/internal/models"
)

type Handler struct {
	userChecker UserChecker
	DB          DBInt
	cacheGetter BannerGetter
}

type RequestBanner struct {
	IsActive  bool         `params:"is_active" json:"is_active"`
	FeatureID uint64       `params:"feature_id" json:"feature_id"`
	TagIDs    []uint64     `params:"tag_ids" json:"tag_ids"`
	Content   pgtype.JSONB `params:"content" json:"content"`
}

func NewHandler(db DBInt, cache BannerGetter, checker UserChecker) *Handler {
	return &Handler{
		DB:          db,
		cacheGetter: cache,
		userChecker: checker,
	}
}

func (h *Handler) GetUserBanner(c *fiber.Ctx) error {
	tagID := c.QueryInt("tag_id")
	featureID := c.QueryInt("feature_id")
	lastRevision := c.QueryBool("use_last_revision", false)

	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}

	data, err := h.cacheGetter.GetUserBanner(uint64(tagID), uint64(featureID), lastRevision, admin)

	if err != nil {
		log.Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Type("json").Send(data)
}

func (h *Handler) PostBanner(c *fiber.Ctx) error {
	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}
	if !admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	banner := RequestBanner{}

	if err = c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return err
	}

	result, err := h.DB.PostBanner(models.Banner{
		IsActive:  banner.IsActive,
		TagIDs:    banner.TagIDs,
		FeatureID: banner.FeatureID,
		Content:   banner.Content,
	})

	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"banner_id": result})
}

func (h *Handler) PatchBanner(c *fiber.Ctx) error {
	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}
	if !admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	banner := RequestBanner{}

	if err = c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := c.ParamsInt("id", -1)

	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.DB.PatchBanner(models.Banner{
		IsActive:  banner.IsActive,
		TagIDs:    banner.TagIDs,
		FeatureID: banner.FeatureID,
		Content:   banner.Content,
	}, id)

	if err != nil {
		log.Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteBanner(c *fiber.Ctx) error {
	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}
	if !admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	id, err := c.ParamsInt("id", -1)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.DB.DeleteBanner(uint64(id))

	if err != nil {
		log.Error(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
