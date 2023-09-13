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

var TEST_DATA_PERSON = []models.Person{
	{
		Name:    "Dave",
		Email:   "dave@email.com",
		Phone:   "555-555-5522",
		Profile: "Dave is a software developer with 5 years of experience.",
	}, {
		Name:    "Dan",
		Email:   "dan@email.com",
		Phone:   "555-555-5533",
		Profile: "Dan is a construction finisher with 10 years of experience.",
	}, {
		Name:    "Drew",
		Email:   "drew@email.com",
		Phone:   "555-555-5544",
		Profile: "Drew is a civil engineer with 15 years of experience.",
	}, {
		Name:    "Dylan",
		Email:   "dylan@email.com",
		Phone:   "555-555-1155",
		Profile: "Dylan is a sailor with 20 years of experience.",
	},
}

// Create a new person
func (suite *TestWithDbSuite) TestCreatePerson() {
	// calculate the total people in the database before adding a new one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(TEST_DATA_PERSON[0])

	req := httptest.NewRequest(http.MethodPost, "/person", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parseResponseData(suite.T(), resp, models.Person{}).(models.Person)

	// calculate the total people in the database after adding a new one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people increased by 1
	assert.Equal(suite.T(), len(totalPeopleBefore)+1, len(totalPeopleAfter))

	// Confirm that the person in the response matches the person in the database
	assert.Equal(suite.T(), respPerson.Name, TEST_DATA_PERSON[0].Name)
}

// Test to get a person by Id
func (suite *TestWithDbSuite) TestReadPerson() {
	// Create a person to read
	personAdded := TEST_DATA_PERSON[1]
	database.DAL.CreatePerson(&personAdded)

	// Create a request to the person route
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parseResponseData(suite.T(), resp, models.Person{}).(models.Person)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), TEST_DATA_PERSON[1].Name, respPerson.Name)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// Test to update a person by Id
func (suite *TestWithDbSuite) TestUpdatePerson() {
	// Create a person to update
	personAdded := TEST_DATA_PERSON[2]
	database.DAL.CreatePerson(&personAdded)

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(TEST_DATA_PERSON[2])

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/person/%v", personAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parseResponseData(suite.T(), resp, models.Person{}).(models.Person)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), TEST_DATA_PERSON[2].Name, respPerson.Name)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestDeletePerson tests the DELETE /person/:id route
func (suite *TestWithDbSuite) TestDeletePerson() {
	// Create a person to delete
	personAdded := TEST_DATA_PERSON[3]
	database.DAL.CreatePerson(&personAdded)

	// calculate the total people in the database before deleting one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// calculate the total people in the database after deleting one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people decreased by 1
	assert.Equal(suite.T(), len(totalPeopleBefore)-1, len(totalPeopleAfter))

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err == nil {
		suite.T().Error(err)
	}
}
