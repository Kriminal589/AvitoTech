package handlers

import (
	"AvitoTech/internal/handlers/contract"
	"AvitoTech/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgtype"
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

	message, err := h.DB.GetUserBanner(uint64(tagId), uint64(featureId))

	if err != nil {
		return err
	}
	return c.JSON(message)
}

func (h *Handler) GetBanner(c *fiber.Ctx) error {
	var message []pgtype.JSONB
	var err error

	tagId := c.QueryInt("tag_id", -1)
	featureId := c.QueryInt("feature_id", -1)

	if tagId == -1 && featureId != -1 {
		message, err = h.DB.GetBanner(uint64(featureId))
	}

	if err != nil {
		return err
	}
	return c.JSON(message)
}

func (h *Handler) PostBanner(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) PatchBanner(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) DeleteBanner(c *fiber.Ctx) error {
	params := models.DeleteBanner{}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	fmt.Println(params.ID)

	return err
}
