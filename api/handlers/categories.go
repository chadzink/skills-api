package handlers

import (
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary Create a new category
// @Description Create a new category entity
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body interface{} true "Category JSON object that needs to be created"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /category [post]
func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return ErorrAndDataResponse(c, err, category)
	}

	database.DAL.CreateCategory(category)

	return DataResponse(c, category)
}

// @Summary Create one or more new categories
// @Description Create one or more new category entities
// @Tags Categories
// @Accept json
// @Produce json
// @Param categories body []interface{} true "Array of Category objects in JSON that need to be created"
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /categories [post]
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

// @Summary List a category by id
// @Description List a category by id
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /category/{id} [get]
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

// @Summary List all categories
// @Description List all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /categories [get]
func ListCategories(c *fiber.Ctx) error {
	// TODO: Add pagination

	categories, err := database.DAL.GetAllCategories()

	if err != nil {
		return ErorrAndDataResponse(c, err, categories)
	}

	return DataResponse(c, categories)
}

// @Summary Update a category by id
// @Description Update a category by id
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body interface{} true "Category JSON object that needs to be updated"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /category/{id} [post]
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

// @Summary Delete a category by id
// @Description Delete a category by id
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /category/{id} [delete]
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
