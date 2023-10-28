package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"skills-api/models"

	"github.com/stretchr/testify/assert"
)

var TEST_DATA_USERS = models.User{
	DisplayName: "John Doe",
	Email:       "john.doe@email.com",
	Password:    "password",
}

var TEST_DATA_LOGIN_REQUESTS = models.LoginRequest{
	Email:    "test.user@email.com",
	Password: "testpassword",
}

// Test register new user request
func (suite *TestWithDbSuite) TestRegisterNewUser() {
	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_USERS)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Check if the new user was added to the database
	var user models.User
	suite.db.Where("email = ?", TEST_DATA_USERS.Email).First(&user)
	assert.Equal(suite.T(), TEST_DATA_USERS.Email, user.Email)

	// Confirm that the response body has a token
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var loginResponse models.LoginResponse
		if err := json.Unmarshal(bodyBytes, &loginResponse); err != nil {
			suite.T().Error(err)
		}

		assert.NotEmpty(suite.T(), loginResponse.Token, "Token should not be empty")
	}
}

// Test login request
func (suite *TestWithDbSuite) TestLogin() {
	// Create a request to the category route
	reqBodyJson, _ := json.Marshal(TEST_DATA_LOGIN_REQUESTS)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Confirm that the response body has a token
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var loginResponse models.LoginResponse
		if err := json.Unmarshal(bodyBytes, &loginResponse); err != nil {
			suite.T().Error(err)
		}

		assert.NotEmpty(suite.T(), loginResponse.Token, "Token should not be empty")
	}
}
