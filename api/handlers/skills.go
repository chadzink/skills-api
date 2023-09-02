package handlers

import (
	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/models"
	"github.com/gofiber/fiber/v2"
)

func CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)

	if err := c.BodyParser(skill); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    skill,
		})
	}

	database.DB.Db.Create(&skill)

	return c.Status(200).JSON(skill)
}

func CreateSkills(c *fiber.Ctx) error {
	parsedSkills := new([]models.Skill)

	if err := c.BodyParser(parsedSkills); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    parsedSkills,
		})
	}

	createdSkills := make([]models.Skill, len(*parsedSkills))

	for index, skill := range *parsedSkills {
		database.DB.Db.Create(&skill)
		createdSkills[index] = skill
	}

	return c.Status(200).JSON(createdSkills)
}

func ListSkill(c *fiber.Ctx) error {
	id := c.Params("id")

	skill, err := database.DB.GetSkillById(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    skill,
		})
	}

	return c.Status(200).JSON(skill)
}

func ListSkills(c *fiber.Ctx) error {
	skills, err := database.DB.GetAllSkills()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    skills,
		})
	}

	return c.Status(200).JSON(skills)
}
