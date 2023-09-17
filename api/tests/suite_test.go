package tests

import (
	"log"
	"testing"

	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create a new test suite
type TestWithDbSuite struct {
	suite.Suite
	app *fiber.App
	db  *gorm.DB
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TestWithDbSuite) SetupTest() {
	// Create a new fiber app
	suite.app = fiber.New()
	setupTestRoutes(suite.app)

	// Open an SQLite in-memory database for testing
	conn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	suite.db = conn
	database.DAL.Db = suite.db
	database.MigrateDb(suite.db)

	seedDatabase(&database.DAL)
}

func (suite *TestWithDbSuite) TearDownSuite() {
	// Close the database connection
	sqlDB, err := suite.db.DB()
	if err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	sqlDB.Close()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestWithDbRunSuite(t *testing.T) {
	suite.Run(t, new(TestWithDbSuite))
}

func setupTestRoutes(a *fiber.App) {

	// Set up the routes for categories
	a.Post("/category", handlers.CreateCategory)
	a.Get("/category/:id", handlers.ListCategory)
	a.Post("/category/:id", handlers.UpdateCategory)
	a.Delete("/category/:id", handlers.DeleteCategory)
	a.Get("/categories", handlers.ListCategories)
	a.Post("/categories", handlers.CreateCategories)

	// Set up the routes for people
	a.Post("/person", handlers.CreatePerson)
	a.Get("/person/:id", handlers.ListPerson)
	a.Post("/person/:id", handlers.UpdatePerson)
	a.Delete("/person/:id", handlers.DeletePerson)

	// Set up the routes for skills
	a.Post("/skill", handlers.CreateSkill)
	a.Get("/skill/:id", handlers.ListSkill)
	a.Post("/skill/:id", handlers.UpdateSkill)
	a.Delete("/skill/:id", handlers.DeleteSkill)
	a.Get("/skills", handlers.ListSkills)
	a.Post("/skills", handlers.CreateSkills)
}

type TestData struct {
	skills     []models.Skill
	categories []models.Category
	people     []models.Person
}

func seedDatabase(dal *database.DataAccessLayer) {

	testData := TestData{
		categories: []models.Category{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:        "Programming Language",
				Description: "A programming language is a formal language comprising a set of instructions that produce various kinds of output. Programming languages are used in computer programming to implement algorithms.",
				ShortKey:    "prog_lang",
				Active:      true,
			}, {
				Model: gorm.Model{
					ID: 2,
				},
				Name:        "Databases",
				Description: "A database is an organized collection of structured information, or data, typically stored electronically in a computer system. A database is usually controlled by a database management system.",
				ShortKey:    "db",
				Active:      true,
			}, {
				Model: gorm.Model{
					ID: 3,
				},
				Name:        "Operating Systems",
				Description: "An operating system is system software that manages computer hardware, software resources, and provides common services for computer programs.",
				ShortKey:    "os",
				Active:      true,
			}, {
				Model: gorm.Model{
					ID: 4,
				},
				Name:        "Cloud Providers",
				Description: "A cloud provider is a company that offers some component of cloud computing -- typically infrastructure as a service, software as a service or platform as a service -- to other businesses or individuals.",
				ShortKey:    "cp",
				Active:      true,
			},
		},
		skills: []models.Skill{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:        "Go",
				Description: "Go is a compiled, statically typed programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
				ShortKey:    "go",
				Active:      true,
				Categories:  []*models.Category{},
			}, {
				Model: gorm.Model{
					ID: 2,
				},
				Name:        "JavaScript",
				Description: "JavaScript, often abbreviated as JS, is a programming language that conforms to the ECMAScript specification. JavaScript is high-level, often just-in-time compiled, and multi-paradigm. It has curly-bracket syntax, dynamic typing, prototype-based object-orientation, and first-class functions.",
				ShortKey:    "js",
				Active:      true,
			}, {
				Model: gorm.Model{
					ID: 3,
				},
				Name:        "Python",
				Description: "Python is an interpreted, high-level and general-purpose programming language. Python's design philosophy emphasizes code readability with its notable use of significant indentation.",
				ShortKey:    "py",
				Active:      true,
			}, {
				Model: gorm.Model{
					ID: 4,
				},
				Name:        "Java",
				Description: "Java is a class-based, object-oriented programming language that is designed to have as few implementation dependencies as possible. It is a general-purpose programming language intended to let application developers write once, run anywhere, meaning that compiled Java code can run on all platforms that support Java without the need for recompilation.",
				ShortKey:    "java",
				Active:      true,
			},
		},
		people: []models.Person{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:    "John",
				Email:   "john@email.com",
				Phone:   "555-555-5555",
				Profile: "John is a software developer with 10 years of experience.",
			}, {
				Model: gorm.Model{
					ID: 2,
				},
				Name:    "Jane",
				Email:   "jane@email.com",
				Phone:   "555-555-5555",
				Profile: "Jane is a software developer with 15 years of experience.",
			}, {
				Model: gorm.Model{
					ID: 3,
				},
				Name:    "Joe",
				Email:   "joe@email.com",
				Phone:   "555-555-5555",
				Profile: "Joe is an IT manager with 20 years of experience.",
			},
		},
	}

	// Link the skills and categories
	// testData.skills[0].Categories = append(testData.skills[0].Categories, &testData.categories[0])

	// Add categories to database
	for _, category := range testData.categories {
		dal.CreateCategory(&category)
	}

	// Add skills to database
	for _, skill := range testData.skills {
		dal.CreateSkill(&skill)
	}

	// Add people to database
	for _, person := range testData.people {
		dal.CreatePerson(&person)
	}
}
