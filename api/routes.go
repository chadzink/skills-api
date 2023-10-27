package main

import (
	"skills-api/auth"
	"skills-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	// Note: This is just an example, please use a secure secret key
	jwt := auth.NewAuthMiddleware(auth.JWTSecretKey)

	app.Get("/", handlers.Default)

	// Create a Login route
	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/register", handlers.RegisterNewUser)

	// CRUD for skill entity
	app.Post("/skill", jwt, handlers.CreateSkill)
	app.Get("/skill/:id", jwt, handlers.ListSkill)
	app.Post("/skill/:id", jwt, handlers.UpdateSkill)
	app.Delete("/skill/:id", jwt, handlers.DeleteSkill)
	// batch features
	app.Get("/skills", jwt, handlers.ListSkills)
	app.Post("/skills", jwt, handlers.CreateSkills)

	// CRUD for category entity
	app.Post("/category", jwt, handlers.CreateCategory)
	app.Get("/category/:id", jwt, handlers.ListCategory)
	app.Post("/category/:id", jwt, handlers.UpdateCategory)
	app.Delete("/category/:id", jwt, handlers.DeleteCategory)
	// batch features
	app.Get("/categories", jwt, handlers.ListCategories)
	app.Post("/categories", jwt, handlers.CreateCategories)

	// CRUD for person entity
	app.Post("/person", jwt, handlers.CreatePerson)
	app.Get("/person/:id", jwt, handlers.ListPerson)
	app.Post("/person/:id", jwt, handlers.UpdatePerson)
	app.Delete("/person/:id", jwt, handlers.DeletePerson)
	// batch features
	app.Get("/people", jwt, handlers.ListPeople)

	// READ for expertise entity
	app.Get("/expertises", jwt, handlers.ListExpertises)
}
