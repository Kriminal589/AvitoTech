package integrationtests

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

const timeExp = time.Hour * 24

func GenerateToken(id uint64) string {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(timeExp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return "Bearer " + s
}
