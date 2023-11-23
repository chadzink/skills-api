package handlers

import (
	"skills-api/database"

	"github.com/gofiber/fiber/v2"
)

// @Summary List all expertises
// @Description List all expertises
// @Tags Expertises
// @Produce json
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /expertises [get]
func ListExpertises(c *fiber.Ctx) error {
	expertises, err := database.DAL.GetAllExpertise()

	if err != nil {
		return ErorrAndDataResponse(c, err, expertises)
	}

	return DataResponse(c, expertises)
}
