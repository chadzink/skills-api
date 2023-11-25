package main

import (
	"os"

	"skills-api/auth"
	"skills-api/database"
	"skills-api/handlers"

	_ "skills-api/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger" // swagger embed files

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	authMiddleware := auth.NewAuthMiddleware(auth.JWTSecretKey)

	app.Get("/", handlers.Default)
	app.Get("/version", handlers.Version)

	// Create a Login route
	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/register", handlers.RegisterNewUser)

	// Set up the routes for create user API key
	app.Post("/user/api_key", authMiddleware, handlers.CreateAPIKey)

	// CRUD for skill entity
	app.Post("/skill", authMiddleware, handlers.CreateSkill)
	app.Get("/skill/:id", authMiddleware, handlers.ListSkill)
	app.Post("/skill/:id", authMiddleware, handlers.UpdateSkill)
	app.Delete("/skill/:id", authMiddleware, handlers.DeleteSkill)
	// batch features
	app.Get("/skills", authMiddleware, handlers.ListSkills)
	app.Post("/skills", authMiddleware, handlers.CreateSkills)

	// CRUD for category entity
	app.Post("/category", authMiddleware, handlers.CreateCategory)
	app.Get("/category/:id", authMiddleware, handlers.ListCategory)
	app.Post("/category/:id", authMiddleware, handlers.UpdateCategory)
	app.Delete("/category/:id", authMiddleware, handlers.DeleteCategory)
	// batch features
	app.Get("/categories", authMiddleware, handlers.ListCategories)
	app.Post("/categories", authMiddleware, handlers.CreateCategories)

	// CRUD for person entity
	app.Post("/person", authMiddleware, handlers.CreatePerson)
	app.Get("/person/:id", authMiddleware, handlers.ListPerson)
	app.Post("/person/:id", authMiddleware, handlers.UpdatePerson)
	app.Delete("/person/:id", authMiddleware, handlers.DeletePerson)
	// batch features
	app.Get("/people", authMiddleware, handlers.ListPeople)

	// READ for expertise entity
	app.Get("/expertises", authMiddleware, handlers.ListExpertises)
}

// @title Skills API Swagger Documentation
// @version 0.1.0
// @description This is an example project built with Fiber and GORM to demonstrate how to build RESTful APIs with Go.
// @termsOfService http://swagger.io/terms/
// @contact.name Chad Zink
// @contact.url http://www.chadzink.com
// @contact.email me@chadzink.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	dbConnectError := database.ConnectDb()
	if dbConnectError != nil {
		os.Exit(2)
	}

	app := fiber.New()

	// Serve the Swagger UI at /swagger/index.html
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	setupRoutes(app)

	app.Listen(":3000")
}
