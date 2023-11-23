package handlers

import "github.com/gofiber/fiber/v2"

// @Summary Home
// @Description Home (default) route
// @Tags Home
// @Produce text/plain
// @Success 200 {object} string
// @Router / [get]
func Default(c *fiber.Ctx) error {
	return c.SendString("Welcome to the skils API!")
}

// @Summary Version
// @Description Version route
// @Tags Home
// @Produce json
// @Success 200 {object} interface{}
// @Router /version [get]
func Version(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"version": "0.1.0",
	})
}
