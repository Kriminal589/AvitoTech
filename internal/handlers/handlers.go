package handlers

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"AvitoTech/internal/models"
)

const timeExp = time.Hour * 24

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

func (h *Handler) Login(c *fiber.Ctx) error {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}

	type user struct {
		ID uint64 `json:"id"`
	}

	var body user
	err := c.BodyParser(&body)

	if err != nil {
		err = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "cannot parse JSON body",
		})
		if err != nil {
			return err
		}

		return err
	}

	claims := jwt.MapClaims{
		"id":  body.ID,
		"exp": time.Now().Add(timeExp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	err = c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
	})
	if err != nil {
		return err
	}

	return err
}
