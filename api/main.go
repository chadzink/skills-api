package main

import (
	"os"

	"github.com/chadzink/skills-api/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	dbConnectError := database.ConnectDb()
	if dbConnectError != nil {
		os.Exit(2)
	}

	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
