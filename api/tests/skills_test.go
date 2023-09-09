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
var testSkillData = map[string]interface{}{
	"name":        "Go",
	"description": "Go is a compiled, statically typed programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
	"short_key":   "go",
	"active":      true,
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

func TestCreateSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Setup skill route
	app.Post("/skill", handlers.CreateSkill)

	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	// Create a request to the skill route
	reqBodyJson, _ := json.Marshal(testSkillData)

	req := httptest.NewRequest(http.MethodPost, "/skill", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseSkillFromResponse(t, resp)

	// Check if the response returned the correct data Name property
	assert.Equal(t, testSkillData["name"], respSkill.Name)

	// Check if the number of skills in the database increased by one
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, len(totalSkillsBefore)+1, len(skills))
	}
}

func TestReadSkill(t *testing.T) {
	// init app and database
	app := SetupTestAppAndDatabase(t)

	// Create a skill to read
	skillAdded := ConvertMapToSkill(testSkillData)
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
	assert.Equal(t, testSkillData["name"], respSkill.Name)
}
