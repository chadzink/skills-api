package tests

import (
	"io"
	"net/http/httptest"

	"github.com/chadzink/skills-api/handlers"
	"github.com/stretchr/testify/assert" // add Testify package
)

func (suite *TestWithDbSuite) TestDefault() {
	suite.app.Get("/", handlers.Default)

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Confirm that the response body is "Welcoem to the skils API!"
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		suite.T().Error(err)
	} else {
		bodyString := string(bodyBytes)
		assert.Equal(suite.T(), "Welcome to the skils API!", bodyString)
	}
}
