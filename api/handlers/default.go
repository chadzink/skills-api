package handlers

import "github.com/gofiber/fiber/v2"

func Default(c *fiber.Ctx) error {
	return c.SendString("Welcome to the skils API!")
}
