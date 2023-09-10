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

var testPersonData = []map[string]interface{}{
	{
		"name":    "John",
		"email":   "john@email.com",
		"phone":   "555-555-5555",
		"profile": "John is a software developer with 10 years of experience.",
	}, {
		"name":    "Jane",
		"email":   "jane@email.com",
		"phone":   "555-555-5555",
		"profile": "Jane is a software developer with 15 years of experience.",
	}, {
		"name":    "Joe",
		"email":   "joe@email.com",
		"phone":   "555-555-5555",
		"profile": "Joe is a software developer with 20 years of experience.",
	},
}

func ConvertMapToPerson(m map[string]interface{}) models.Person {
	return models.Person{
		Name:    m["name"].(string),
		Email:   m["email"].(string),
		Phone:   m["phone"].(string),
		Profile: m["profile"].(string),
	}
}

func parsePersonFromResponse(t *testing.T, resp *http.Response) models.Person {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult

		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// Extract the response data into a person object
		responseDataAsPerson := ConvertMapToPerson(responseObject.Data)
		responseDataAsPerson.ID = uint(responseObject.Data["ID"].(float64))

		return responseDataAsPerson
	}

	return models.Person{}
}

// Create a new person
func TestCreatePerson(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup person route
	app.Post("/person", handlers.CreatePerson)

	// calculate the total people in the database before adding a new one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(testPersonData[0])

	req := httptest.NewRequest(http.MethodPost, "/person", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parsePersonFromResponse(t, resp)

	// calculate the total people in the database after adding a new one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people increased by 1
	assert.Equal(t, len(totalPeopleBefore)+1, len(totalPeopleAfter))

	// Confirm that the person in the response matches the person in the database
	assert.Equal(t, respPerson.Name, totalPeopleAfter[len(totalPeopleAfter)-1].Name)
}

// Test to get a person by Id
func TestReadPerson(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup person route
	app.Get("/person/:id", handlers.ListPerson)

	// Create a person to read
	personAdded := ConvertMapToPerson(testPersonData[0])
	database.DAL.CreatePerson(&personAdded)

	// Create a request to the person route
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parsePersonFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testPersonData[0]["name"], respPerson.Name)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		t.Error(err)
	}
}

// Test to update a person by Id
func TestUpdatePerson(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup person route
	app.Post("/person/:id", handlers.UpdatePerson)

	// Create a person to update
	personAdded := ConvertMapToPerson(testPersonData[0])
	database.DAL.CreatePerson(&personAdded)

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(testPersonData[1])

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/person/%v", personAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the person from the response
	respPerson := parsePersonFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testPersonData[1]["name"], respPerson.Name)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		t.Error(err)
	}
}

// TestDeletePerson tests the DELETE /person/:id route
func TestDeletePerson(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup person route
	app.Delete("/person/:id", handlers.DeletePerson)

	// Create a person to delete
	personAdded := ConvertMapToPerson(testPersonData[0])
	database.DAL.CreatePerson(&personAdded)

	// calculate the total people in the database before deleting one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// calculate the total people in the database after deleting one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people decreased by 1
	assert.Equal(t, len(totalPeopleBefore)-1, len(totalPeopleAfter))

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err == nil {
		t.Error(err)
	}
}
