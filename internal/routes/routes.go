package routes

import (
	"github.com/gofiber/fiber/v2"

	"AvitoTech/internal/handlers"
	"AvitoTech/internal/middleware"
)

func InitializeRoutes(a *fiber.App, handler *handlers.Handler) {
	route := a.Group("/api")
	jwtMW := middleware.JWTProtected()

	route.Get("/user_banner", jwtMW, handler.GetUserBanner)
	route.Get("/banner", jwtMW, handler.GetBanners)
	route.Post("/banner", jwtMW, handler.PostBanner)
	route.Patch("/banner/:id", jwtMW, handler.PatchBanner)
	route.Delete("/banner/:id", jwtMW, handler.DeleteBanner)
	route.Post("/login", handler.Login)
}
