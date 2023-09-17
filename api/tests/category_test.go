package tests

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"

	"net/http"
	"net/http/httptest"

	"github.com/chadzink/skills-api/database"
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
		SkillIds:    []uint{2, 3},
	},
	{
		Name:        "Operations",
		Description: "Operations is the process of operating a business.",
		ShortKey:    "operations",
		Active:      true,
		SkillIds:    []uint{3, 4},
	},
	{
		Name:        "Communication",
		Description: "Communication is the process of communicating with others.",
		ShortKey:    "communication",
		Active:      true,
		SkillIds:    []uint{1, 4},
	},
}

// Create a new category
func (suite *TestWithDbSuite) TestCreateCategory() {
	// suite.updateGoldenFile = true

	// calculate the total categories in the database before adding a new one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_CATEGORIES[0])

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// calculate the total categories in the database after adding a new one
	totalCategoriesAfter, _ := database.DAL.GetAllCategories()

	// Confirm that the total categories increased by 1
	assert.Equal(suite.T(), len(totalCategoriesBefore)+1, len(totalCategoriesAfter))

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Create Category", "create_category", resp)
}

// Test to get a category by Id
func (suite *TestWithDbSuite) TestReadCategory() {
	// suite.updateGoldenFile = true

	// Create a category to read
	categoryAdded := TEST_DATA_CATEGORIES[1]
	database.DAL.CreateCategory(&categoryAdded)

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		suite.T().Error(err)
	}

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Read Category", "read_category", resp)
}

// Test to update a category by Id
func (suite *TestWithDbSuite) TestUpdateCategory() {
	// suite.updateGoldenFile = true

	// Create a category to update
	categoryAdded := TEST_DATA_CATEGORIES[2]
	database.DAL.CreateCategory(&categoryAdded)

	// Change the name of the category
	categoryAdded.Name = "Updated Category Name"

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(categoryAdded)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/category/%v", categoryAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Update Category", "update_category", resp)

	// Check if the category was found in the database
	if _, err := database.DAL.GetCategoryById(categoryAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestDeleteCategory tests the DELETE /category/:id route
func (suite *TestWithDbSuite) TestDeleteCategory() {
	// suite.updateGoldenFile = true

	// Create a category to delete
	categoryAdded := TEST_DATA_CATEGORIES[3]
	database.DAL.CreateCategory(&categoryAdded)

	// calculate the total categories in the database before deleting one
	totalCategoriesBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/category/%v", categoryAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Delete Category", "delete_category", resp)

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
	// suite.updateGoldenFile = true

	// calculate the total categories in the database before adding a new one
	totalCategories, _ := database.DAL.GetAllCategories()

	// Create a category to read if none exist
	if len(totalCategories) == 0 {
		for _, categoryAdded := range TEST_DATA_CATEGORIES {
			database.DAL.CreateCategory(&categoryAdded)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponsesToGoldenFile("Test Read Categories", "read_categories", resp)
}

// TestCreateCategories tests the POST /categories route
func (suite *TestWithDbSuite) TestCreateCategories() {
	// suite.updateGoldenFile = true

	// calculate the total categories in the database before adding a new one
	totalCategorysBefore, _ := database.DAL.GetAllCategories()

	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_SKILLS)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponsesToGoldenFile("Test Create Categories", "create_categories", resp)

	// Check if the number of categories in the database increased by the total number of categories added
	if categories, err := database.DAL.GetAllCategories(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalCategorysBefore)+len(TEST_DATA_SKILLS), len(categories))
	}
}
