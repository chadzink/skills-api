package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
	"github.com/stretchr/testify/assert" // add Testify package
)

// Test data for skills
var testSkillData = []map[string]interface{}{
	{
		"name":        "Go",
		"description": "Go is a compiled, statically typed programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
		"short_key":   "go",
		"active":      true,
	}, {
		"name":        "JavaScript",
		"description": "JavaScript, often abbreviated as JS, is a programming language that conforms to the ECMAScript specification. JavaScript is high-level, often just-in-time compiled, and multi-paradigm. It has curly-bracket syntax, dynamic typing, prototype-based object-orientation, and first-class functions.",
		"short_key":   "js",
		"active":      true,
	}, {
		"name":        "Python",
		"description": "Python is an interpreted, high-level and general-purpose programming language. Python's design philosophy emphasizes code readability with its notable use of significant indentation.",
		"short_key":   "py",
		"active":      true,
	}, {
		"name":        "Java",
		"description": "Java is a class-based, object-oriented programming language that is designed to have as few implementation dependencies as possible. It is a general-purpose programming language intended to let application developers write once, run anywhere, meaning that compiled Java code can run on all platforms that support Java without the need for recompilation.",
		"short_key":   "java",
		"active":      true,
	},
}

func ConvertMapToSkill(m map[string]interface{}) models.Skill {
	return models.Skill{
		Name:        m["name"].(string),
		Description: m["description"].(string),
		ShortKey:    m["short_key"].(string),
		Active:      m["active"].(bool),
	}
}

func parseSkillFromResponse(t *testing.T, resp *http.Response) models.Skill {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult

		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// Extract the response data into a skill object
		responseDataAsSkill := ConvertMapToSkill(responseObject.Data)
		responseDataAsSkill.ID = uint(responseObject.Data["ID"].(float64))

		return responseDataAsSkill
	}

	return models.Skill{}
}

// TestCreateSkill tests the POST /skill route
func TestCreateSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup skill route
	app.Post("/skill", handlers.CreateSkill)

	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	// Create a request to the skill route
	reqBodyJson, _ := json.Marshal(testSkillData[0])

	req := httptest.NewRequest(http.MethodPost, "/skill", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseSkillFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testSkillData[0]["name"], respSkill.Name)

	// Check if the number of skills in the database increased by one
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, len(totalSkillsBefore)+1, len(skills))
	}
}

// TestListSkill tests the GET /skill/:id route
func TestReadSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a skill to read
	skillAdded := ConvertMapToSkill(testSkillData[0])
	database.DAL.CreateSkill(&skillAdded)

	// Setup skill route
	app.Get("/skill/:id", handlers.ListSkill)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseSkillFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testSkillData[0]["name"], respSkill.Name)
}

// TestUpdateSkill tests the POST /skill/:id route
func TestUpdateSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a skill to update
	skillAdded := ConvertMapToSkill(testSkillData[0])
	database.DAL.CreateSkill(&skillAdded)

	// Setup skill route
	app.Post("/skill/:id", handlers.UpdateSkill)

	// Update the skill
	skillAdded.Description = "Updated description"
	reqBodyJson, _ := json.Marshal(skillAdded)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/skill/%v", skillAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseSkillFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, skillAdded.Description, respSkill.Description)
}

// TestDeleteSkill tests the DELETE /skill/:id route
func TestDeleteSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a skill to delete
	skillAdded := ConvertMapToSkill(testSkillData[0])
	database.DAL.CreateSkill(&skillAdded)

	// Setup skill route
	app.Delete("/skill/:id", handlers.DeleteSkill)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Check if the skill was deleted
	if _, err := database.DAL.GetSkillById(skillAdded.ID); err == nil {
		t.Error("Expected error when getting deleted skill")
	} else {
		assert.Equal(t, "record not found", err.Error())
	}
}

// TestListSkills tests the GET /skills route
func TestListSkills(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a skill to read
	for _, skillData := range testSkillData {
		skillAdded := ConvertMapToSkill(skillData)
		database.DAL.CreateSkill(&skillAdded)
	}

	// Setup skill route
	app.Get("/skills", handlers.ListSkills)

	req := httptest.NewRequest(http.MethodGet, "/skills", nil)
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skill from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a skill object
			responseDataAsSkill := ConvertMapToSkill(responseObject)
			responseDataAsSkill.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(t, testSkillData[i]["name"], responseDataAsSkill.Name)
		}
	}
}

// TestCreateSkills tests the POST /skills route
func TestCreateSkills(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup skill route
	app.Post("/skills", handlers.CreateSkills)

	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	// Create a request to the skill route
	reqBodyJson, _ := json.Marshal(testSkillData)

	req := httptest.NewRequest(http.MethodPost, "/skills", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skills from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a skill object
			responseDataAsSkill := ConvertMapToSkill(responseObject)
			responseDataAsSkill.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(t, testSkillData[i]["name"], responseDataAsSkill.Name)
		}
	}

	// Check if the number of skills in the database increased by the total number of skills added
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, len(totalSkillsBefore)+len(testSkillData), len(skills))
	}
}
