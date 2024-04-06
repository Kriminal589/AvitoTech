package routes

import (
	"AvitoTech/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(a *fiber.App, handler *handlers.Handler) {
	route := a.Group("/api")

	route.Get("/user_banner", handler.GetUserBanner)
	route.Get("/banner", handler.GetBanner)
	route.Post("/banner", handler.PostBanner)
	route.Patch("/banner/:id", handler.PatchBanner)
	route.Delete("/banner/:id", handler.DeleteBanner)
}
