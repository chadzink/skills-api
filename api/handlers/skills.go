package handlers

import (
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateSkill creates a new skill entity
func CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)

	if err := c.BodyParser(skill); err != nil {
		return ErorrAndDataResponse(c, err, skill)
	}

	database.DAL.CreateSkill(skill)

	return DataResponse(c, skill)
}

// CreateSkills creates one or more new skill entities
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

// ListSkill lists a skill by id
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

// ListSkills lists all skills
func ListSkills(c *fiber.Ctx) error {
	// TO DO: Add paging to this endpoint

	skills, err := database.DAL.GetAllSkills()

	if err != nil {
		return ErorrAndDataResponse(c, err, skills)
	}

	return DataResponse(c, skills)
}

// UpdateSkill updates a skill by id
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

// DeleteSkill deletes a skill by id
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
