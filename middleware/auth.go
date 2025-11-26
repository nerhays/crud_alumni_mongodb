package middleware

import (
	"crud_alumni/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Token diperlukan"})
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{"error": "Format token salah"})
        }

        claims, err := utils.ValidateToken(tokenParts[1])
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Token invalid"})
        }

        c.Locals("user_id", claims.UserID)
        c.Locals("username", claims.Username)
        c.Locals("role", claims.Role)
        return c.Next()
    }
}

func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("role").(string)
        if role != "admin" {
            return c.Status(403).JSON(fiber.Map{"error": "Hanya admin"})
        }
        return c.Next()
    }
}
