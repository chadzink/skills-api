package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"skills-api/database"
	"skills-api/models"

	"github.com/stretchr/testify/assert" // add Testify package
)

var TEST_DATA_SKILLS = []models.Skill{
	{
		Name:        "Boat Building",
		Description: "Boat building is the design and construction of boats and their systems. This includes at a minimum a hull, with propulsion, mechanical, navigation, safety and other systems as a craft requires.",
		ShortKey:    "boat-building",
		Active:      true,
		CategoryIds: []uint{1, 2},
	},
	{
		Name:        "Drywall Mudding",
		Description: "Drywall mudding is the process of applying drywall mud to drywall seams.",
		ShortKey:    "drywall-mudding",
		Active:      true,
		CategoryIds: []uint{2},
	},
	{
		Name:        "Roadmap Planning",
		Description: "Roadmap planning is the process of planning a roadmap.",
		ShortKey:    "roadmap-planning",
		Active:      true,
		CategoryIds: []uint{3, 4},
	},
	{
		Name:        "Data Modeling",
		Description: "Data modeling is the process of creating a data model for the data to be stored in a database.",
		ShortKey:    "data-modeling",
		Active:      true,
		CategoryIds: []uint{3},
	},
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *TestWithDbSuite) TestCreateSkill() {
	// suite.updateGoldenFile = true

	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	skillToAdd := TEST_DATA_SKILLS[0]

	//convert skill to add into javascript object for ody of request
	reqBodyJson, _ := json.Marshal(skillToAdd)

	req := suite.GetJwtRequest(http.MethodPost, "/skill", bytes.NewReader(reqBodyJson))
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Create Skill", "create_skill", resp)

	// Check if the number of skills in the database increased by one
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalSkillsBefore)+1, len(skills))
	}
}

// TestListSkill tests the GET /skill/:id route
func (suite *TestWithDbSuite) TestReadSkill() {
	// suite.updateGoldenFile = true

	// Create a skill to read
	skillAdded := TEST_DATA_SKILLS[1]

	database.DAL.CreateSkill(&skillAdded)

	req := suite.GetJwtRequest(http.MethodGet, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Read Skill", "read_skill", resp)

	// Check if the skill was found in the database
	if _, err := database.DAL.GetSkillById(skillAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestUpdateSkill tests the POST /skill/:id route
func (suite *TestWithDbSuite) TestUpdateSkill() {
	// suite.updateGoldenFile = true

	// Create a skill to update
	skillAdded := TEST_DATA_SKILLS[2]

	database.DAL.CreateSkill(&skillAdded)

	// Change the name, description, & the categories of the skill
	skillAdded.Name = "Updated Name"
	skillAdded.Categories = nil
	skillAdded.CategoryIds = []uint{1, 2, 3}
	skillAdded.Description = "Updated description"

	reqBodyJson, _ := json.Marshal(skillAdded)
	req := suite.GetJwtRequest(http.MethodPost, fmt.Sprintf("/skill/%v", skillAdded.ID), bytes.NewReader(reqBodyJson))
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Update Skill", "update_skill", resp)
}

// TestDeleteSkill tests the DELETE /skill/:id route
func (suite *TestWithDbSuite) TestDeleteSkill() {
	// suite.updateGoldenFile = true

	// Create a skill to delete
	skillAdded := TEST_DATA_SKILLS[3]

	database.DAL.CreateSkill(&skillAdded)

	req := suite.GetJwtRequest(http.MethodDelete, fmt.Sprintf("/skill/%v", skillAdded.ID), nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Delete Skill", "delete_skill", resp)

	// Check if the skill was deleted
	if _, err := database.DAL.GetSkillById(skillAdded.ID); err == nil {
		suite.T().Error("Expected error when getting deleted skill")
	} else {
		assert.Equal(suite.T(), "record not found", err.Error())
	}
}

// TestListSkills tests the GET /skills route
func (suite *TestWithDbSuite) TestListSkills() {
	// suite.updateGoldenFile = true

	// calculate the total skills in the database before adding a new one
	totalSkills, _ := database.DAL.GetAllSkills()

	// Create a skill to read if none exist
	if len(totalSkills) == 0 {
		for _, skillAdded := range TEST_DATA_SKILLS {
			database.DAL.CreateSkill(&skillAdded)
		}
	}

	req := suite.GetJwtRequest(http.MethodGet, "/skills", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponsesToGoldenFile("Test List Skills", "list_skills", resp)
}

// TestCreateSkills tests the POST /skills route
func (suite *TestWithDbSuite) TestCreateSkills() {
	// suite.updateGoldenFile = true

	// calculate the total skills in the database before adding a new one
	totalSkillsBefore, _ := database.DAL.GetAllSkills()

	// Create a request to the skill route
	reqBodyJson, _ := json.Marshal(TEST_DATA_SKILLS)

	req := suite.GetJwtRequest(http.MethodPost, "/skills", bytes.NewReader(reqBodyJson))
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponsesToGoldenFile("Test Create Skills", "create_skills", resp)

	// Check if the number of skills in the database increased by the total number of skills added
	if skills, err := database.DAL.GetAllSkills(); err != nil {
		suite.T().Error(err)
	} else {
		assert.Equal(suite.T(), len(totalSkillsBefore)+len(TEST_DATA_SKILLS), len(skills))
	}
}
