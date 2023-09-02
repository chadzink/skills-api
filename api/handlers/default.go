package handlers

import "github.com/gofiber/fiber/v2"

func Default(c *fiber.Ctx) error {
	return c.SendString("Welcoem to the skils API!")
}
