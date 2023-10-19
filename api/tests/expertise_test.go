package tests

import (
	"net/http"

	"github.com/stretchr/testify/assert" // add Testify package
)

// TestListExpertise tests the GET /skills route
func (suite *TestWithDbSuite) TestListExpertise() {
	suite.updateGoldenFile = true

	// Create a request to the expertise route
	req := suite.GetJwtRequest(http.MethodGet, "/expertises", nil)
	resp, _ := suite.app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Comapre the response to the golden file
	suite.CheckResponsesToGoldenFile("Test List Expertises", "list_expertises", resp)
}
