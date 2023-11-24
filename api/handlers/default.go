package handlers

import "github.com/gofiber/fiber/v2"

// @Summary Home
// @Description Home (default) route
// @Tags Home
// @Produce text/plain
// @Success 200 {string} string "Welcome to the skils API!"
// @Router / [get]
func Default(c *fiber.Ctx) error {
	return c.SendString("Welcome to the skils API!")
}

type VersionResponse struct {
	Version string `json:"version" example:"0.1.0"`
}

// @Summary Version
// @Description Version route
// @Tags Home
// @Produce json
// @Success 200 {object} VersionResponse
// @Router /version [get]
func Version(c *fiber.Ctx) error {
	return c.Status(200).JSON(VersionResponse{
		Version: "0.1.0",
	})
}
