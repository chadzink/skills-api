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

	// CRUD for person entity
	app.Post("/person", handlers.CreatePerson)
	app.Get("/person/:id", handlers.ListPerson)
	app.Post("/person/:id", handlers.UpdatePerson)
	app.Delete("/person/:id", handlers.DeletePerson)
	// batch features
	app.Get("/people", handlers.ListPeople)

	// READ for expertise entity
	app.Get("/expertises", handlers.ListExpertises)
}
