package handlers

import (
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory creates a new category entity
func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return ErorrAndDataResponse(c, err, category)
	}

	database.DAL.CreateCategory(category)

	return DataResponse(c, category)
}

// CreateCategories creates one or more new category entities
func CreateCategories(c *fiber.Ctx) error {
	parsedCategories := new([]models.Category)

	if err := c.BodyParser(parsedCategories); err != nil {
		return ErorrAndDataResponse(c, err, parsedCategories)
	}

	createdCategories := make([]models.Category, len(*parsedCategories))

	for index, category := range *parsedCategories {
		database.DAL.CreateCategory(&category)
		createdCategories[index] = category
	}

	return DataResponse(c, createdCategories)
}

// ListCategory lists a category by id
func ListCategory(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	category, err := database.DAL.GetCategoryById(id)

	if err != nil {
		return ErorrAndDataResponse(c, err, category)
	}

	return DataResponse(c, category)
}

// ListCategories lists all categories
func ListCategories(c *fiber.Ctx) error {
	// TODO: Add pagination

	categories, err := database.DAL.GetAllCategories()

	if err != nil {
		return ErorrAndDataResponse(c, err, categories)
	}

	return DataResponse(c, categories)
}

// UpdateCategory updates a category by id
func UpdateCategory(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return ErorrAndDataResponse(c, err, category)
	}

	err = database.DAL.UpdateCategoryById(id, category)

	if err != nil {
		return ErorrAndDataResponse(c, err, category)
	}

	return DataResponse(c, category)
}

// DeleteCategory deletes a category by id
func DeleteCategory(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	err = database.DAL.DeleteCategoryById(id)

	if err != nil {
		return ErorrAndDataResponse(c, err, id)
	}

	return DeletedResponse(c, "Category deleted", id)
}
