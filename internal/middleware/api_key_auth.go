package middleware

import (
	"sompong-api/internal/config"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth verifies X-API-Key header against QUIZ_API_KEY.
func APIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" || apiKey != config.App.QuizAPIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid api key",
			})
		}
		return c.Next()
	}
}
