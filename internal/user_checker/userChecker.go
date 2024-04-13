package user_checker

import (
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

func (u UserChecker) IsAdmin(c *fiber.Ctx) (bool, error) {
	user := c.Locals(u.contextKey).(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	admin, err := u.userDB.GetUserRole(uint64(claims["id"].(float64)))

	return admin, err
}
