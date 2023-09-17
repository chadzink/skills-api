package tests

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
	"github.com/stretchr/testify/assert"
	// add Testify package
)

var TEST_DATA_CATEGORIES = []models.Category{
	{
		Name:        "Construction",
		Description: "Construction is the process of constructing a building or infrastructure.",
		ShortKey:    "construction",
		Active:      true,
		SkillIds:    []uint{1, 2},
	},
	{
		Name:        "Leadership",
		Description: "Leadership is the ability to lead others.",
		ShortKey:    "leadership",
		Active:      true,
	},
	{
		Name:        "Operations",
		Description: "Operations is the process of operating a business.",
		ShortKey:    "operations",
		Active:      true,
	},
	{
		Name:        "Communication",
		Description: "Communication is the process of communicating with others.",
		ShortKey:    "communication",
		Active:      true,
	},
}

// Create a new category
func (suite *TestWithDbSuite) TestCreateCategory() {
	// calculate the total categories in the database before adding a new one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_CATEGORIES[0])

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseResponseData(suite.T(), resp, models.Category{}).(models.Category)

	// calculate the total categories in the database after adding a new one
	totalCategoriesAfter, _ := database.DAL.GetAllCategories()

	// Confirm that the total categories increased by 1
	assert.Equal(suite.T(), len(totalCategoriesBefore)+1, len(totalCategoriesAfter))

	// Confirm that the category in the response matches the category in the database
	assert.Equal(suite.T(), respCategory.Name, TEST_DATA_CATEGORIES[0].Name)

	// Check if the new category in database is linked to the skills in the id list
	if category, err := database.DAL.GetCategoryById(respCategory.ID); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), category.Skills[0].ID, TEST_DATA_CATEGORIES[0].SkillIds[0])
		assert.Equal(suite.T(), category.Skills[1].ID, TEST_DATA_CATEGORIES[0].SkillIds[1])
	}
}

// Test to get a category by Id
func (suite *TestWithDbSuite) TestReadCategory() {
	// Create a category to read
	categoryAdded := TEST_DATA_CATEGORIES[1]
	database.DAL.CreateCategory(&categoryAdded)

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseResponseData(suite.T(), resp, models.Category{}).(models.Category)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), TEST_DATA_CATEGORIES[1].Name, respCategory.Name)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// Test to update a category by Id
func (suite *TestWithDbSuite) TestUpdateCategory() {
	// Create a category to update
	categoryAdded := TEST_DATA_CATEGORIES[2]
	database.DAL.CreateCategory(&categoryAdded)

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_CATEGORIES[2])

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/category/%v", categoryAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the category from the response
	respCategory := parseResponseData(suite.T(), resp, models.Category{}).(models.Category)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), TEST_DATA_CATEGORIES[2].Name, respCategory.Name)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestDeleteCategory tests the DELETE /category/:id route
func (suite *TestWithDbSuite) TestDeleteCategory() {
	// Create a category to delete
	categoryAdded := TEST_DATA_CATEGORIES[3]
	database.DAL.CreateCategory(&categoryAdded)

	// calculate the total categories in the database before deleting one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// calculate the total categories in the database after deleting one
	totalCategoriesAfter, _ := database.DAL.GetAllCategories()

	// Confirm that the total categories decreased by 1
	assert.Equal(suite.T(), len(totalCategoriesBefore)-1, len(totalCategoriesAfter))

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err == nil {
		suite.T().Error(err)
	}
}

// TestListCategories tests the GET /categories route
func (suite *TestWithDbSuite) TestListCategories() {
	// calculate the total categories in the database before adding a new one
	totalCategories, _ := database.DAL.GetAllCategories()

	// Create a category to read if none exist
	if len(totalCategories) == 0 {
		for _, categoryAdded := range TEST_DATA_CATEGORIES {
			database.DAL.CreateCategory(&categoryAdded)
		}

		totalCategories, _ = database.DAL.GetAllCategories()
	}

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the category from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			suite.T().Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a category object
			responseDataAsCategory := MapToCategory(responseObject)
			responseDataAsCategory.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(suite.T(), totalCategories[i].Name, responseDataAsCategory.Name)
		}
	}
}

// TestCreateCategories tests the POST /categories route
func (suite *TestWithDbSuite) TestCreateCategories() {
	// calculate the total categories in the database before adding a new one
	totalCategorysBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_SKILLS)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the categories from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			suite.T().Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a category object
			responseDataAsCategory := MapToCategory(responseObject)
			responseDataAsCategory.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(suite.T(), TEST_DATA_SKILLS[i].Name, responseDataAsCategory.Name)
		}
	}

	// Check if the number of categories in the database increased by the total number of categories added
	if categories, err := database.DAL.GetAllCategories(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalCategorysBefore)+len(TEST_DATA_SKILLS), len(categories))
	}
}
