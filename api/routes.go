package main

import (
	"github.com/chadzink/skills-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Default)

	// CRUD for skill entity
	app.Post("/skill", handlers.CreateSkill)
	app.Get("/skill/:id", handlers.ListSkill)
	app.Post("/skill/:id", handlers.UpdateSkill)
	app.Delete("/skill/:id", handlers.DeleteSkill)
	// batch features
	app.Get("/skills", handlers.ListSkills)
	app.Post("/skills", handlers.CreateSkills)

	// CRUD for category entity
	app.Post("/category", handlers.CreateCategory)
	app.Get("/category/:id", handlers.ListCategory)
	app.Post("/category/:id", handlers.UpdateCategory)
	app.Delete("/category/:id", handlers.DeleteCategory)
	// batch features
	app.Get("/categories", handlers.ListCategories)
	app.Post("/categories", handlers.CreateCategories)
}
