package handlers

import (
	"AvitoTech/internal/handlers/contract"
	"AvitoTech/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"time"
)

type Handler struct {
	DB contract.DBInt
}

func NewHandler(DB contract.DBInt) *Handler {
	return &Handler{DB: DB}
}

func (h *Handler) GetUserBanner(c *fiber.Ctx) error {
	tagId := c.QueryInt("tag_id")
	featureId := c.QueryInt("feature_id")

	data, err := h.DB.GetUserBanner(uint64(tagId), uint64(featureId))

	if err != nil {
		return err
	}
	return c.JSON(data)
}

func (h *Handler) GetBanners(c *fiber.Ctx) error {
	var data []models.Banner
	var err error

	tagId := c.QueryInt("tag_id", -1)
	featureId := c.QueryInt("feature_id", -1)

	if tagId == -1 && featureId != -1 {
		data, err = h.DB.GetBannersWithFeatureId(uint64(featureId))
	} else if tagId != -1 && featureId == -1 {

	}

	if err != nil {
		return err
	}
	return c.JSON(data)
}

func (h *Handler) PostBanner(c *fiber.Ctx) error {
	banner := models.RequestBanner{}

	if err := c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return err
	}

	result, err := h.DB.PostBanner(banner)

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"banner_id": result})
}

func (h *Handler) PatchBanner(c *fiber.Ctx) error {
	banner := models.RequestBanner{}

	if err := c.BodyParser(&banner); err != nil {
		log.Errorf("Unable to parse request body: %v", err)
		return err
	}

	id, err := c.ParamsInt("id", -1)

	if err != nil {
		return err
	}

	result, err := h.DB.PatchBanner(banner, id)

	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *Handler) DeleteBanner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", -1)

	if err != nil {
		return err
	}

	result, err := h.DB.DeleteBanner(uint64(id))

	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}

	type user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		"email": body.Email,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString([]byte("AvitoTech"))
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

func (h *Handler) Restricted(c *fiber.Ctx) error {
	err := c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
	if err != nil {
		return err
	}

	return nil
}
