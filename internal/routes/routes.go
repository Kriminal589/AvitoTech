package routes

import (
	"AvitoTech/internal/handlers"
	"AvitoTech/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(a *fiber.App, handler *handlers.Handler) {
	route := a.Group("/api")

	route.Get("/user_banner", handler.GetUserBanner)
	route.Get("/banner", handler.GetBanners)
	route.Post("/banner", middleware.JWTProtected(), handler.PostBanner)
	route.Patch("/banner/:id", middleware.JWTProtected(), handler.PatchBanner)
	route.Delete("/banner/:id", middleware.JWTProtected(), handler.DeleteBanner)
	route.Post("/login", handler.Login)
}
