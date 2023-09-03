package handlers

import (
	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/models"
	"github.com/gofiber/fiber/v2"
)

// CreateCategory creates a new category entity
func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    category,
		})
	}

	database.DAL.CreateCategory(category)

	return c.Status(200).JSON(category)
}

// CreateCategories creates one or more new category entities
func CreateCategories(c *fiber.Ctx) error {
	parsedCategories := new([]models.Category)

	if err := c.BodyParser(parsedCategories); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    parsedCategories,
		})
	}

	createdCategories := make([]models.Category, len(*parsedCategories))

	for index, category := range *parsedCategories {
		database.DAL.CreateCategory(&category)
		createdCategories[index] = category
	}

	return c.Status(200).JSON(createdCategories)
}

// ListCategory lists a category by id
func ListCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	category, err := database.DAL.GetCategoryById(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    category,
		})
	}

	return c.Status(200).JSON(category)
}

// ListCategories lists all categories
func ListCategories(c *fiber.Ctx) error {
	// TODO: Add pagination

	categories, err := database.DAL.GetAllCategories()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    categories,
		})
	}

	return c.Status(200).JSON(categories)
}

// UpdateCategory updates a category by id
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    category,
		})
	}

	err := database.DAL.UpdateCategoryById(id, category)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    category,
		})
	}

	return c.Status(200).JSON(category)
}

// DeleteCategory deletes a category by id
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	err := database.DAL.DeleteCategoryById(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    id,
		})
	}

	return c.Status(200).JSON(id)
}
