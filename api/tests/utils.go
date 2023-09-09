package tests

import (
	"testing"

	"github.com/chadzink/skills-api/database"
	"github.com/gofiber/fiber/v2"
)

func SetupTestAppAndDatabase(t *testing.T) *fiber.App {
	// Create a new fiber app and test database connection
	app := fiber.New()
	dbConnectError := database.ConnectTestDb()

	if dbConnectError != nil {
		t.Error(dbConnectError)
	}

	return app
}
