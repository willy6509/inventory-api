package middleware

import (
	"inventory-api/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Token Required"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	secret := config.GetEnv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Token Invalid"})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_role", claims["role"]) // Simpan role untuk dipakai di controller
	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("user_role")
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden - Admin Access Only"})
	}
	return c.Next()
}