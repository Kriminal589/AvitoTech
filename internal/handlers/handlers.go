package handlers

import (
	"AvitoTech/internal/databases"
	"AvitoTech/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	DB databases.DBInt
}

func NewHandler(DB databases.DBInt) *Handler {
	return &Handler{DB: DB}
}

func (h *Handler) GetUserBanner(c *fiber.Ctx) error {
	params := models.UserBanner{}

	err := c.QueryParser(&params)
	bannerId, err := h.DB.GetUserBanner(params.TagID, params.FeatureID)

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"banner_id": bannerId,
	})
}

func (h *Handler) GetBanner(c *fiber.Ctx) error {
	return nil
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
