package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v2"
)

func JWTProtected() func(c *fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey:   []byte("AvitoTech"),
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
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
