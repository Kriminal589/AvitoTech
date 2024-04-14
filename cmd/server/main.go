package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"AvitoTech/internal/cache"
	"AvitoTech/internal/databases"
	"AvitoTech/internal/handlers"
	"AvitoTech/internal/handlers/getbanner"
	"AvitoTech/internal/middleware"
	checker "AvitoTech/internal/user-checker"
)

func main() {
	app := fiber.New()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"))

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	pgxDB := databases.NewPgxDB(pool, zapadapter.NewLogger(zap.NewNop()))

	bannerCache := cache.New(pgxDB)
	userChecker := checker.New(pgxDB, os.Getenv("CONTEXT_JWT_KEY"))

	handler := handlers.NewHandler(pgxDB, bannerCache, userChecker)
	getBannerHandler := getbanner.NewHandler(pgxDB, userChecker)

	route := app.Group("/api")
	jwtMW := middleware.JWTProtected()

	route.Get("/user_banner", jwtMW, handler.GetUserBanner)
	route.Get("/banner", jwtMW, getBannerHandler.GetBanners)
	route.Post("/banner", jwtMW, handler.PostBanner)
	route.Patch("/banner/:id", jwtMW, handler.PatchBanner)
	route.Delete("/banner/:id", jwtMW, handler.DeleteBanner)

	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
