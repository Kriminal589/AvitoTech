package userchecker

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserChecker struct {
	userDB     UserDB
	contextKey string
}

func New(userDB UserDB, contextKey string) *UserChecker {
	return &UserChecker{userDB, contextKey}
}

var errCastingFailure = errors.New("casting error")

func (u UserChecker) IsAdmin(c *fiber.Ctx) (bool, error) {
	user, ok := c.Locals(u.contextKey).(*jwt.Token)

	if !ok {
		return false, errCastingFailure
	}

	claims, ok := user.Claims.(jwt.MapClaims)

	if !ok {
		return false, errCastingFailure
	}

	admin, err := u.userDB.GetUserRole(uint64(claims["id"].(float64)))

	return admin, err
}
