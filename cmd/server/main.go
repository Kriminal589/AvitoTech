package main

import (
	"AvitoTech/internal/databases"
	"AvitoTech/internal/handlers"
	"AvitoTech/internal/routes"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"log"
	"os"
)

var embedMigrations embed.FS

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

	var db *sql.DB
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, os.Getenv("PATH_MIGRATIONS")); err != nil {
		panic(err)
	}

	handler := handlers.NewHandler(pgxDB)
	routes.InitializeRoutes(app, handler)

	if err := app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
