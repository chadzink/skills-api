package tests

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/chadzink/skills-api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert" // add Testify package
)

func TestDefault(t *testing.T) {
	app := fiber.New()

	app.Get("/", handlers.Default)

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)

	// Confirm that the response status code is 200
	assert.Equal(t, 200, resp.StatusCode)

	// Confirm that the response body is "Welcoem to the skils API!"
	defer resp.Body.Close()
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		bodyString := string(bodyBytes)
		assert.Equal(t, "Welcome to the skils API!", bodyString)
	}
}
