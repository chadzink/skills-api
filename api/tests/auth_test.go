package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"skills-api/handlers"
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

// Test create API key request, this is a protected route that requires a valid JWT token for specific user
func (suite *TestWithDbSuite) TestCreateAPIKeyForUser() {
	apiKeyRequestForUser := models.NewAPIKeyRequest{
		Email:    TEST_DATA_LOGIN_REQUESTS.Email,
		Password: TEST_DATA_LOGIN_REQUESTS.Password,
		// Expires on 3 days from current date and time
		ExpiresOn: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+3, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Now().Location()),
	}

	// Create a Login request
	reqBodyJson, _ := json.Marshal(apiKeyRequestForUser)

	// Create a request to the user API key route
	req := suite.GetJwtRequest(http.MethodPost, "/user/api_key", bytes.NewReader(reqBodyJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Check if the response fetched a new user api key fron the database
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var respResult handlers.ResponseResult[models.UserAPIKey]
		if err := json.Unmarshal(bodyBytes, &respResult); err != nil {
			suite.T().Error(err)
		}

		var userApiKey models.UserAPIKey = respResult.Data

		assert.NotEmpty(suite.T(), userApiKey.Key, "API Key should not be empty")

		// Check if a request can be made for skills using the new API key
		req = httptest.NewRequest(http.MethodGet, "/skills", nil)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("X-API-Key", userApiKey.Key)
		req.Header.Set("X-API-Email", TEST_DATA_LOGIN_REQUESTS.Email)
		resp, _ = suite.app.Test(req)

		// Confirm that the response status code is 200
		assert.Equal(suite.T(), 200, resp.StatusCode)
	}
}

// Test get current user by token and api key
func (suite *TestWithDbSuite) TestGetCurrentUser() {
	// Create a Login request
	reqBodyJson, _ := json.Marshal(TEST_DATA_LOGIN_REQUESTS)
	var loginResponse models.LoginResponse

	// Create a request to the login route
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
		if err := json.Unmarshal(bodyBytes, &loginResponse); err != nil {
			suite.T().Error(err)
		}

		assert.NotEmpty(suite.T(), loginResponse.Token, "Token should not be empty")
	}

	// Create a request to the user route
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)
	resp, _ = suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Confirm that the response body has a token
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		var respResult handlers.ResponseResult[models.User]
		if err := json.Unmarshal(bodyBytes, &respResult); err != nil {
			suite.T().Error(err)
		}

		var user models.User = respResult.Data

		assert.NotEmpty(suite.T(), user.Email, "Email should not be empty")
	}
}
