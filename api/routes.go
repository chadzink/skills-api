package main

import (
	"github.com/chadzink/skills-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Default)

	app.Post("/skill", handlers.CreateSkill)
	app.Get("/skill/:id", handlers.ListSkill)
	app.Get("/skills", handlers.ListSkills)
	app.Post("/skills", handlers.CreateSkills)

}
