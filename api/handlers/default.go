package handlers

import "github.com/gofiber/fiber/v2"

func Default(c *fiber.Ctx) error {
	return c.SendString("Welcome to the skils API!")
}

func Version(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"version": "0.1.0",
	})
}
