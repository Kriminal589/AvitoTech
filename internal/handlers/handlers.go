package handlers

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"AvitoTech/internal/models"
)

type Handler struct {
	userChecker UserChecker
	Db          DBInt
	cacheGetter BannerGetter
}

func NewHandler(DB DBInt, cache BannerGetter, checker UserChecker) *Handler {
	return &Handler{
		Db:          DB,
		cacheGetter: cache,
		userChecker: checker,
	}
}

func (h *Handler) GetUserBanner(c *fiber.Ctx) error {
	tagId := c.QueryInt("tag_id")
	featureId := c.QueryInt("feature_id")
	lastRevision := c.QueryBool("use_last_revision", false)

	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}

	data, err := h.cacheGetter.GetUserBanner(uint64(tagId), uint64(featureId), lastRevision, admin)

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

	tagId := c.QueryInt("tag_id", -1)
	featureId := c.QueryInt("feature_id", -1)
	limit := c.QueryInt("limit", -1)
	offset := c.QueryInt("offset", 0)

	admin, err := h.userChecker.IsAdmin(c)

	if err != nil {
		return err
	}
	if !admin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if tagId == -1 && featureId != -1 {
		data, err = h.Db.GetBannersByFeatureId(uint64(featureId), limit, offset)
	}
	if tagId != -1 && featureId == -1 {
		data, err = h.Db.GetBannersByTagId(uint64(tagId), limit, offset)
	}
	if tagId != -1 && featureId != -1 {
		data, err = h.Db.GetBanners(uint64(featureId), uint64(tagId), limit, offset)
	}
	if tagId == -1 && featureId == -1 {
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

	banner := models.RequestBanner{}

	if err := c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return err
	}

	result, err := h.Db.PostBanner(banner)

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

	banner := models.RequestBanner{}

	if err := c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := c.ParamsInt("id", -1)

	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.Db.PatchBanner(banner, id)

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

	err = h.Db.DeleteBanner(uint64(id))

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
		Id uint64 `json:"id"`
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
		"id":  body.Id,
		"exp": time.Now().Add(time.Hour * 3).Unix(),
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
