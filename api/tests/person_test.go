package tests

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"

	"net/http"

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
		PersonSkills: []*models.PersonSkill{
			{
				SkillID:     1,
				ExpertiseID: 5,
			},
			{
				SkillID:     2,
				ExpertiseID: 4,
			},
		},
	}, {
		Name:    "Dan",
		Email:   "dan@email.com",
		Phone:   "555-555-5533",
		Profile: "Dan is a construction finisher with 10 years of experience.",
		PersonSkills: []*models.PersonSkill{
			{
				SkillID:     2,
				ExpertiseID: 3,
			},
			{
				SkillID:     3,
				ExpertiseID: 2,
			},
		},
	}, {
		Name:    "Drew",
		Email:   "drew@email.com",
		Phone:   "555-555-5544",
		Profile: "Drew is a civil engineer with 15 years of experience.",
		PersonSkills: []*models.PersonSkill{
			{
				SkillID:     3,
				ExpertiseID: 1,
			},
			{
				SkillID:     4,
				ExpertiseID: 5,
			},
		},
	}, {
		Name:    "Dylan",
		Email:   "dylan@email.com",
		Phone:   "555-555-1155",
		Profile: "Dylan is a sailor with 20 years of experience.",
		PersonSkills: []*models.PersonSkill{
			{
				SkillID:     1,
				ExpertiseID: 1,
			},
			{
				SkillID:     2,
				ExpertiseID: 2,
			},
			{
				SkillID:     3,
				ExpertiseID: 3,
			},
			{
				SkillID:     4,
				ExpertiseID: 4,
			},
		},
	},
}

// Create a new person
func (suite *TestWithDbSuite) TestCreatePerson() {
	// suite.updateGoldenFile = true

	// calculate the total people in the database before adding a new one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(TEST_DATA_PERSON[0])

	req := suite.GetJwtRequest(http.MethodPost, "/person", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Create Person", "create_person", resp)

	// calculate the total people in the database after adding a new one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people increased by 1
	assert.Equal(suite.T(), len(totalPeopleBefore)+1, len(totalPeopleAfter))
}

// Test to get a person by Id
func (suite *TestWithDbSuite) TestReadPerson() {
	// suite.updateGoldenFile = true

	// Create a person to read
	personAdded := TEST_DATA_PERSON[1]
	database.DAL.CreatePerson(&personAdded)

	// Create a request to the person route
	req := suite.GetJwtRequest(http.MethodGet, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Read Person", "read_person", resp)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// Test to update a person by Id
func (suite *TestWithDbSuite) TestUpdatePerson() {
	// suite.updateGoldenFile = true

	// Create a person to update
	personAdded := TEST_DATA_PERSON[2]
	database.DAL.CreatePerson(&personAdded)

	// Change the person's name, email, and add a new skill
	personAdded.Name = "Richard"
	personAdded.Email = "richard@email.com"
	personAdded.PersonSkills = append(personAdded.PersonSkills, &models.PersonSkill{
		SkillID:     1,
		ExpertiseID: 1,
	})

	// Create a request to the person route
	reqBodyJson, _ := json.Marshal(personAdded)

	req := suite.GetJwtRequest(http.MethodPost, fmt.Sprintf("/person/%v", personAdded.ID), bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Update Person", "update_person", resp)

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err != nil {
		suite.T().Error(err)
	}
}

// TestDeletePerson tests the DELETE /person/:id route
func (suite *TestWithDbSuite) TestDeletePerson() {
	// suite.updateGoldenFile = true

	// Create a person to delete
	personAdded := TEST_DATA_PERSON[3]
	database.DAL.CreatePerson(&personAdded)

	// calculate the total people in the database before deleting one
	totalPeopleBefore, _ := database.DAL.GetAllPeople()

	// Create a request to the person route
	req := suite.GetJwtRequest(http.MethodDelete, fmt.Sprintf("/person/%v", personAdded.ID), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponseToGoldenFile("Test Delete Person", "delete_person", resp)

	// calculate the total people in the database after deleting one
	totalPeopleAfter, _ := database.DAL.GetAllPeople()

	// Confirm that the total people decreased by 1
	assert.Equal(suite.T(), len(totalPeopleBefore)-1, len(totalPeopleAfter))

	// Check if the person was found in the database
	if _, err := database.DAL.GetPersonById(personAdded.ID); err == nil {
		suite.T().Error(err)
	}
}
