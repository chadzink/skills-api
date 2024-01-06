package middleware

import (
	"skills-api/auth"
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func NewBearerAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
	})
}

func NewAPIKeyAuthMiddleware(apiEmail, apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Setup verify user API request object
		requestVerifyApiKey := models.VerifyAPIKeyRequest{
			Email: apiEmail,
			Key:   apiKey,
		}

		userApiKey, err := database.DAL.VerifyUserAPIKey(requestVerifyApiKey)
		if err != nil && userApiKey.Key != "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid API Key",
			})
		}

		return c.Next()
	}
}

// Middleware JWT function
func NewAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiEmail := c.Get("X-API-Email")

		if apiEmail != "" {
			apiKey := c.Get("X-API-Key")
			return NewAPIKeyAuthMiddleware(apiEmail, apiKey)(c)
		} else {
			return NewBearerAuthMiddleware(auth.JWTSecretKey)(c)
		}
	}
}
