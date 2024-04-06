package main

import (
	"AvitoTech/internal/handlers"
	"AvitoTech/internal/routes"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	app := fiber.New()

	handler := &handlers.Handler{}

	routes.InitializeRoutes(app, handler)

	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
