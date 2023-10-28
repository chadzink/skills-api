package handlers

import (
	"skills-api/database"

	"github.com/gofiber/fiber/v2"
)

// ListExpertises lists all expertise
func ListExpertises(c *fiber.Ctx) error {
	expertises, err := database.DAL.GetAllExpertise()

	if err != nil {
		return ErorrAndDataResponse(c, err, expertises)
	}

	return DataResponse(c, expertises)
}
