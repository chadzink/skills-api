package tests

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
	"github.com/stretchr/testify/assert" // add Testify package
)

var testCategoryData = []map[string]interface{}{
	{
		"name":        "Programming Languages",
		"description": "A programming language is a formal language comprising a set of instructions that produce various kinds of output. Programming languages are used in computer programming to implement algorithms.",
		"short_key":   "prog_lang",
		"active":      true,
	}, {
		"name":        "Databases",
		"description": "A database is an organized collection of structured information, or data, typically stored electronically in a computer system. A database is usually controlled by a database management system.",
		"short_key":   "db",
		"active":      true,
	}, {
		"name":        "Operating Systems",
		"description": "An operating system is system software that manages computer hardware, software resources, and provides common services for computer programs.",
		"short_key":   "os",
		"active":      true,
	}, {
		"name":        "Web Development",
		"description": "Web development is the work involved in developing a Web site for the Internet or an intranet. Web development can range from developing a simple single static page of plain text to complex Web-based Internet applications, electronic businesses, and social network services.",
		"short_key":   "web_dev",
		"active":      true,
	},
}

func ConvertMapToCategory(m map[string]interface{}) models.Category {
	return models.Category{
		Name:        m["name"].(string),
		Description: m["description"].(string),
		ShortKey:    m["short_key"].(string),
		Active:      m["active"].(bool),
	}
}

func parseCategoryFromResponse(t *testing.T, resp *http.Response) models.Category {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult

		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// Extract the response data into a category object
		responseDataAsCategory := ConvertMapToCategory(responseObject.Data)
		responseDataAsCategory.ID = uint(responseObject.Data["ID"].(float64))

		return responseDataAsCategory
	}

	return models.Category{}
}

// Create a new category
func TestCreateCategory(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup category route
	app.Post("/category", handlers.CreateCategory)

	// calculate the total categories in the database before adding a new one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(testCategoryData[0])

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseCategoryFromResponse(t, resp)

	// calculate the total categories in the database after adding a new one
	totalCategoriesAfter, _ := database.DAL.GetAllCategories()

	// Confirm that the total categories increased by 1
	assert.Equal(t, len(totalCategoriesBefore)+1, len(totalCategoriesAfter))

	// Confirm that the category in the response matches the category in the database
	assert.Equal(t, respCategory.Name, totalCategoriesAfter[len(totalCategoriesAfter)-1].Name)
}

// Test to get a category by Id
func TestReadCategory(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup category route
	app.Get("/category/:id", handlers.ListCategory)

	// Create a category to read
	categoryAdded := ConvertMapToCategory(testCategoryData[0])
	database.DAL.CreateCategory(&categoryAdded)

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseCategoryFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testCategoryData[0]["name"], respCategory.Name)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		t.Error(err)
	}
}

// Test to update a category by Id
func TestUpdateCategory(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup category route
	app.Post("/category/:id", handlers.UpdateCategory)

	// Create a category to update
	categoryAdded := ConvertMapToCategory(testCategoryData[0])
	database.DAL.CreateCategory(&categoryAdded)

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(testCategoryData[1])

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/category/%v", categoryAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseCategoryFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testCategoryData[1]["name"], respCategory.Name)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		t.Error(err)
	}
}

// TestDeleteCategory tests the DELETE /category/:id route
func TestDeleteCategory(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup category route
	app.Delete("/category/:id", handlers.DeleteCategory)

	// Create a category to delete
	categoryAdded := ConvertMapToCategory(testCategoryData[0])
	database.DAL.CreateCategory(&categoryAdded)

	// calculate the total categories in the database before deleting one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// calculate the total categories in the database after deleting one
	totalCategoriesAfter, _ := database.DAL.GetAllCategories()

	// Confirm that the total categories decreased by 1
	assert.Equal(t, len(totalCategoriesBefore)-1, len(totalCategoriesAfter))

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err == nil {
		t.Error(err)
	}
}

// TestListCategories tests the GET /categories route
func TestListCategories(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a category to read
	for _, categoryData := range testCategoryData {
		categoryAdded := ConvertMapToCategory(categoryData)
		database.DAL.CreateCategory(&categoryAdded)
	}

	// Setup category route
	app.Get("/categories", handlers.ListCategories)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the category from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a category object
			responseDataAsCategory := ConvertMapToCategory(responseObject)
			responseDataAsCategory.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(t, testCategoryData[i]["name"], responseDataAsCategory.Name)
		}
	}
}

// TestCreateCategories tests the POST /categories route
func TestCreateCategories(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup category route
	app.Post("/categories", handlers.CreateCategories)

	// calculate the total categories in the database before adding a new one
	totalCategorysBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(testCategoryData)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the categories from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a category object
			responseDataAsCategory := ConvertMapToCategory(responseObject)
			responseDataAsCategory.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(t, testCategoryData[i]["name"], responseDataAsCategory.Name)
		}
	}

	// Check if the number of categories in the database increased by the total number of categories added
	if categories, err := database.DAL.GetAllCategories(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, len(totalCategorysBefore)+len(testCategoryData), len(categories))
	}
}
