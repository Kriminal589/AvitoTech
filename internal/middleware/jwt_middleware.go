package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"os"
)

func JWTProtected() func(c *fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: jwtError,
		ContextKey:   os.Getenv("CONTEXT_JWT_KEY"),
	}

	return jwtware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if errors.Is(err, errors.New("missing or malformed JWT")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
