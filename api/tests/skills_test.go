package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
	"github.com/stretchr/testify/assert" // add Testify package
)

var TEST_DATA_SKILLS = []models.Skill{
	{
		Name:        "Boat Building",
		Description: "Boat building is the design and construction of boats and their systems. This includes at a minimum a hull, with propulsion, mechanical, navigation, safety and other systems as a craft requires.",
		ShortKey:    "boat-building",
		Active:      true,
	},
	{
		Name:        "Drywall Mudding",
		Description: "Drywall mudding is the process of applying drywall mud to drywall seams.",
		ShortKey:    "drywall-mudding",
		Active:      true,
	},
	{
		Name:        "Roadmap Planning",
		Description: "Roadmap planning is the process of planning a roadmap.",
		ShortKey:    "roadmap-planning",
		Active:      true,
	},
	{
		Name:        "Data Modeling",
		Description: "Data modeling is the process of creating a data model for the data to be stored in a database.",
		ShortKey:    "data-modeling",
		Active:      true,
	},
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *TestWithDbSuite) TestCreateSkill() {
	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	skillToAdd := TEST_DATA_SKILLS[0]

	//convert skill to add into javascript object for ody of request
	reqBodyJson, _ := json.Marshal(skillToAdd)

	req := httptest.NewRequest(http.MethodPost, "/skill", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseResponseData(suite.T(), resp, models.Skill{}).(models.Skill)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), skillToAdd.Name, respSkill.Name)

	// Check if the number of skills in the database increased by one
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalSkillsBefore)+1, len(skills))
	}
}

// TestListSkill tests the GET /skill/:id route
func (suite *TestWithDbSuite) TestReadSkill() {
	// Create a skill to read
	skillAdded := TEST_DATA_SKILLS[1]

	database.DAL.CreateSkill(&skillAdded)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseResponseData(suite.T(), resp, models.Skill{}).(models.Skill)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), skillAdded.Name, respSkill.Name)

	// Check if the skill was found in the database
	if _, err := database.DAL.GetSkillById(skillAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestUpdateSkill tests the POST /skill/:id route
func (suite *TestWithDbSuite) TestUpdateSkill() {
	// Create a skill to update
	skillAdded := TEST_DATA_SKILLS[2]

	database.DAL.CreateSkill(&skillAdded)

	// Update the skill
	skillAdded.Description = "Updated description"
	reqBodyJson, _ := json.Marshal(skillAdded)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/skill/%v", skillAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the skill from the response
	respSkill := parseResponseData(suite.T(), resp, models.Skill{}).(models.Skill)

	// Check if the response returned the correct data Name property
	assert.Equal(suite.T(), skillAdded.Description, respSkill.Description)
}

// TestDeleteSkill tests the DELETE /skill/:id route
func (suite *TestWithDbSuite) TestDeleteSkill() {
	// Create a skill to delete
	skillAdded := TEST_DATA_SKILLS[3]

	database.DAL.CreateSkill(&skillAdded)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Check if the skill was deleted
	if _, err := database.DAL.GetSkillById(skillAdded.ID); err == nil {
		suite.T().Error("Expected error when getting deleted skill")
	} else {
		assert.Equal(suite.T(), "record not found", err.Error())
	}
}

// TestListSkills tests the GET /skills route
func (suite *TestWithDbSuite) TestListSkills() {
	// calculate the total skills in the database before adding a new one
	totalSkills, _ := database.DAL.GetAllSkills()

	// Create a skill to read if none exist
	if len(totalSkills) == 0 {
		for _, skillAdded := range TEST_DATA_SKILLS {
			database.DAL.CreateSkill(&skillAdded)
		}

		totalSkills, _ = database.DAL.GetAllSkills()
	}

	req := httptest.NewRequest(http.MethodGet, "/skills", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the skill from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			suite.T().Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a skill object
			responseDataAsSkill := MapToSkill(responseObject)
			responseDataAsSkill.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(suite.T(), totalSkills[i].Name, responseDataAsSkill.Name)
		}
	}
}

// TestCreateSkills tests the POST /skills route
func (suite *TestWithDbSuite) TestCreateSkills() {
	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	// Create a request to the skill route
	reqBodyJson, _ := json.Marshal(TEST_DATA_SKILLS)

	req := httptest.NewRequest(http.MethodPost, "/skills", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Get the skills from the response
	var responseObjects handlers.ResponseResults
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		// parse the response body
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			suite.T().Error(err)
		}

		for i, responseObject := range responseObjects.Data {
			// Extract the response data into a skill object
			responseDataAsSkill := MapToSkill(responseObject)
			responseDataAsSkill.ID = uint(responseObject["ID"].(float64))

			// Check if the response returned the correct data Name property
			assert.Equal(suite.T(), TEST_DATA_SKILLS[i].Name, responseDataAsSkill.Name)
		}
	}

	// Check if the number of skills in the database increased by the total number of skills added
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalSkillsBefore)+len(TEST_DATA_SKILLS), len(skills))
	}
}
