package handlers

import (
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary Create a new skill
// @Description Create a new skill entity
// @Tags Skills
// @Accept json
// @Produce json
// @Param skill body interface{} true "Skill JSON object that needs to be created"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skill [post]
func CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)

	if err := c.BodyParser(skill); err != nil {
		return ErorrAndDataResponse(c, err, skill)
	}

	database.DAL.CreateSkill(skill)

	return DataResponse(c, skill)
}

// @Summary Create one or more new skills
// @Description Create one or more new skill entities
// @Tags Skills
// @Accept json
// @Produce json
// @Param skills body []interface{} true "Array of Skill objects in JSON that need to be created"
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skills [post]
func CreateSkills(c *fiber.Ctx) error {
	parsedSkills := new([]models.Skill)

	if err := c.BodyParser(parsedSkills); err != nil {
		return ErorrAndDataResponse(c, err, parsedSkills)
	}

	createdSkills := make([]models.Skill, len(*parsedSkills))

	for index, skill := range *parsedSkills {
		database.DAL.CreateSkill(&skill)
		createdSkills[index] = skill
	}

	return DataResponse(c, createdSkills)
}

// @Summary List a skill by id
// @Description List a skill by id
// @Tags Skills
// @Accept json
// @Produce json
// @Param id path int true "Skill JSON object"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skill/{id} [get]
func ListSkill(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	skill, err := database.DAL.GetSkillById(id)

	if err != nil {
		return ErorrAndDataResponse(c, err, skill)
	}

	return DataResponse(c, skill)
}

// @Summary List all skills
// @Description List all skills
// @Tags Skills
// @Accept json
// @Produce json
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skills [get]
func ListSkills(c *fiber.Ctx) error {
	// TO DO: Add paging to this endpoint

	skills, err := database.DAL.GetAllSkills()

	if err != nil {
		return ErorrAndDataResponse(c, err, skills)
	}

	return DataResponse(c, skills)
}

// @Summary Update a skill by id
// @Description Update a skill by id
// @Tags Skills
// @Accept json
// @Produce json
// @Param id path int true "Skill ID"
// @Param skill body interface{} true "Skill JSON object that needs to be updated"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skill/{id} [post]
func UpdateSkill(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	skill := new(models.Skill)

	if err := c.BodyParser(skill); err != nil {
		return ErorrAndDataResponse(c, err, skill)
	}

	if err := database.DAL.UpdateSkillById(id, skill); err != nil {
		return ErorrAndDataResponse(c, err, skill)
	}

	// Assign the id to the skill
	skill.ID = id

	return DataResponse(c, skill)
}

// @Summary Delete a skill by id
// @Description Delete a skill by id
// @Tags Skills
// @Accept json
// @Produce json
// @Param id path int true "Skill JSON object"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /skill/{id} [delete]
func DeleteSkill(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	if err := database.DAL.DeleteSkillById(id); err != nil {
		return ErorrAndDataResponse(c, err, id)
	}

	return DeletedResponse(c, "Skill deleted", id)
}
