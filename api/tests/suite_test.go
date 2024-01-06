package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"skills-api/database"
	"skills-api/handlers"
	"skills-api/middleware"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Create a new test suite
type TestWithDbSuite struct {
	suite.Suite
	app              *fiber.App
	db               *gorm.DB
	updateGoldenFile bool
	jwt              string
	apiKey           models.VerifyAPIKeyRequest
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

	suite.jwt = ""
	suite.apiKey = models.VerifyAPIKeyRequest{}

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
// START DEBUG HERE
func TestWithDbRunSuite(t *testing.T) {
	suite.Run(t, new(TestWithDbSuite))
}

func setupTestRoutes(a *fiber.App) {
	authMiddleware := middleware.NewAuthMiddleware()

	// Create a Login route
	a.Post("/auth/login", handlers.Login)
	a.Post("/auth/register", handlers.RegisterNewUser)

	// Set up the routes for create user API key
	a.Get("/user", authMiddleware, handlers.GetCurrentUser)
	a.Post("/user/api_key", authMiddleware, handlers.CreateAPIKey)

	// Set up the routes for categories
	a.Post("/category", authMiddleware, handlers.CreateCategory)
	a.Get("/category/:id", authMiddleware, handlers.ListCategory)
	a.Post("/category/:id", authMiddleware, handlers.UpdateCategory)
	a.Delete("/category/:id", authMiddleware, handlers.DeleteCategory)
	a.Get("/categories", authMiddleware, handlers.ListCategories)
	a.Post("/categories", authMiddleware, handlers.CreateCategories)

	// Set up the routes for people
	a.Post("/person", authMiddleware, handlers.CreatePerson)
	a.Get("/person/:id", authMiddleware, handlers.ListPerson)
	a.Post("/person/:id", authMiddleware, handlers.UpdatePerson)
	a.Delete("/person/:id", authMiddleware, handlers.DeletePerson)

	// Set up the routes for skills
	a.Post("/skill", authMiddleware, handlers.CreateSkill)
	a.Get("/skill/:id", authMiddleware, handlers.ListSkill)
	a.Post("/skill/:id", authMiddleware, handlers.UpdateSkill)
	a.Delete("/skill/:id", authMiddleware, handlers.DeleteSkill)
	a.Get("/skills", authMiddleware, handlers.ListSkills)
	a.Post("/skills", authMiddleware, handlers.CreateSkills)

	// READ for expertise entity
	a.Get("/expertises", authMiddleware, handlers.ListExpertises)
}

func (suite *TestWithDbSuite) CheckResponseToGoldenFile(name string, filename string, resp *http.Response) {
	g := goldie.New(suite.T())

	suite.T().Run(name, func(t *testing.T) {
		// Run your function or code for scenario 1
		result := GetResponseBodyString(t, resp)

		if suite.updateGoldenFile {
			// Use g.Assert to compare the result against the golden file
			g.Update(t, filename, []byte(result))
		} else {
			// Use g.Assert to compare the result against the golden file
			g.Assert(t, filename, []byte(result))
		}
	})
}

func (suite *TestWithDbSuite) CheckResponsesToGoldenFile(name string, filename string, resp *http.Response) {
	g := goldie.New(suite.T())

	suite.T().Run(name, func(t *testing.T) {
		// Run your function or code for scenario 1
		result := GetResponsesBodyString(t, resp)

		if suite.updateGoldenFile {
			// Use g.Assert to compare the result against the golden file
			g.Update(t, filename, []byte(result))
		} else {
			// Use g.Assert to compare the result against the golden file
			g.Assert(t, filename, []byte(result))
		}
	})
}

// This func of the suite is used to login and get a JWT token
func (suite *TestWithDbSuite) GetJwtToken() string {
	// Check if the suite jwt is already set
	if suite.jwt != "" {
		return suite.jwt
	}

	loginInfo := models.LoginRequest{
		Email:    "test.user@email.com",
		Password: "testpassword",
	}

	// Create a Login request
	reqBodyJson, _ := json.Marshal(loginInfo)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Confirm that the response body has a token
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var loginResponse models.LoginResponse
		if err := json.Unmarshal(bodyBytes, &loginResponse); err != nil {
			suite.T().Error(err)
		}

		suite.jwt = loginResponse.Token
	}

	return suite.jwt
}

// This func of the suite is used to create a new user API key and set the verify API key request
func (suite *TestWithDbSuite) GetAPIKey() models.VerifyAPIKeyRequest {
	// Check if the suite jwt is already set
	if suite.apiKey != (models.VerifyAPIKeyRequest{}) {
		return suite.apiKey
	}

	apiKeyRequestForUser := models.NewAPIKeyRequest{
		Email:    "test.user@email.com",
		Password: "testpassword",
		// Expires on 3 days from current date and time
		ExpiresOn: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+3, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Now().Location()),
	}

	// Create a Login request
	reqBodyJson, _ := json.Marshal(apiKeyRequestForUser)

	// Create a request to the user API key route
	req := httptest.NewRequest(http.MethodPost, "/user/api_key", bytes.NewReader(reqBodyJson))
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Confirm that the response body has a token
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var apiKeyResponse models.VerifyAPIKeyRequest
		if err := json.Unmarshal(bodyBytes, &apiKeyResponse); err != nil {
			suite.T().Error(err)
		}

		suite.apiKey = apiKeyResponse
	}

	return suite.apiKey
}

// get an HTTP request object that has the JWT token in the header when given a method, route, and body
func (suit *TestWithDbSuite) GetJwtRequest(method string, route string, body io.Reader) *http.Request {

	req := httptest.NewRequest(method, route, body)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+suit.GetJwtToken())
	return req
}

// get an HTTP request object that has the User API Key in the header when given a method, route, and body
func (suit *TestWithDbSuite) GetUserAPIKeyRequest(method string, route string, body io.Reader) *http.Request {
	// First create or get the user api key
	apiKey := suit.GetAPIKey()

	req := httptest.NewRequest(method, route, body)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-API-Email", apiKey.Email)
	req.Header.Set("X-API-Key", apiKey.Key)
	return req
}

type TestData struct {
	skills     []models.Skill
	categories []models.Category
	people     []models.Person
	users      []models.User
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
				CategoryIds: []uint{1, 2},
			}, {
				Model: gorm.Model{
					ID: 2,
				},
				Name:        "JavaScript",
				Description: "JavaScript, often abbreviated as JS, is a programming language that conforms to the ECMAScript specification. JavaScript is high-level, often just-in-time compiled, and multi-paradigm. It has curly-bracket syntax, dynamic typing, prototype-based object-orientation, and first-class functions.",
				ShortKey:    "js",
				Active:      true,
				CategoryIds: []uint{2, 3},
			}, {
				Model: gorm.Model{
					ID: 3,
				},
				Name:        "Python",
				Description: "Python is an interpreted, high-level and general-purpose programming language. Python's design philosophy emphasizes code readability with its notable use of significant indentation.",
				ShortKey:    "py",
				Active:      true,
				CategoryIds: []uint{3, 4},
			}, {
				Model: gorm.Model{
					ID: 4,
				},
				Name:        "Java",
				Description: "Java is a class-based, object-oriented programming language that is designed to have as few implementation dependencies as possible. It is a general-purpose programming language intended to let application developers write once, run anywhere, meaning that compiled Java code can run on all platforms that support Java without the need for recompilation.",
				ShortKey:    "java",
				Active:      true,
				CategoryIds: []uint{2, 4},
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
				PersonSkills: []*models.PersonSkill{
					{
						SkillID:     1,
						ExpertiseID: 1,
					},
					{
						SkillID:     2,
						ExpertiseID: 2,
					},
				},
			}, {
				Model: gorm.Model{
					ID: 2,
				},
				Name:    "Jane",
				Email:   "jane@email.com",
				Phone:   "555-555-5555",
				Profile: "Jane is a software developer with 15 years of experience.",
				PersonSkills: []*models.PersonSkill{
					{
						SkillID:     2,
						ExpertiseID: 3,
					},
				},
			}, {
				Model: gorm.Model{
					ID: 3,
				},
				Name:    "Joe",
				Email:   "joe@email.com",
				Phone:   "555-555-5555",
				Profile: "Joe is an IT manager with 20 years of experience.",
				PersonSkills: []*models.PersonSkill{
					{
						SkillID:     3,
						ExpertiseID: 4,
					},
					{
						SkillID:     4,
						ExpertiseID: 2,
					},
				},
			},
		},
		users: []models.User{
			{
				Model: gorm.Model{
					ID: 1,
				},
				DisplayName: "Test User",
				Email:       "test.user@email.com",
				Password:    "testpassword",
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

	// Setup users
	for _, user := range testData.users {
		dal.RegisterNewUser(user.Email, user.DisplayName, user.Password)
	}
}
