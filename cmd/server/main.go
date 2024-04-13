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
	"AvitoTech/internal/routes"
	checker "AvitoTech/internal/user_checker"
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
	pgxDB := databases.NewPgxDB(pool, zapadapter.NewLogger(zap.NewNop()))

	bannerCache := cache.New(pgxDB)
	userChecker := checker.New(pgxDB, os.Getenv("CONTEXT_JWT_KEY"))

	handler := handlers.NewHandler(pgxDB, bannerCache, userChecker)
	routes.InitializeRoutes(app, handler)

	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
