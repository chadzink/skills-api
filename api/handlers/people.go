package handlers

import (
	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary Create a new person
// @Description Create a new person entity
// @Tags People
// @Accept json
// @Produce json
// @Param person body interface{} true "Person JSON object that needs to be created"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /person [post]
func CreatePerson(c *fiber.Ctx) error {
	person := new(models.Person)

	if err := c.BodyParser(person); err != nil {
		return ErorrAndDataResponse(c, err, person)
	}

	database.DAL.CreatePerson(person)

	return DataResponse(c, person)
}

// @Summary List a person by id
// @Description List a person by id
// @Tags People
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /person/{id} [get]
func ListPerson(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	person, err := database.DAL.GetPersonById(id)

	if err != nil {
		return ErorrAndDataResponse(c, err, person)
	}

	return DataResponse(c, person)
}

// @Summary List all people
// @Description List all people
// @Tags People
// @Produce json
// @Success 200 {object} []interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /people [get]
func ListPeople(c *fiber.Ctx) error {
	// TODO: Add pagination

	people, err := database.DAL.GetAllPeople()

	if err != nil {
		return ErorrAndDataResponse(c, err, people)
	}

	return DataResponse(c, people)
}

// @Summary Update a person by id
// @Description Update a person by id
// @Tags People
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body interface{} true "Person JSON object that needs to be updated"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /person/{id} [put]
func UpdatePerson(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	person := new(models.Person)

	if err := c.BodyParser(person); err != nil {
		return ErorrAndDataResponse(c, err, person)
	}

	err = database.DAL.UpdatePersonById(id, person)

	if err != nil {
		return ErorrAndDataResponse(c, err, person)
	}

	return DataResponse(c, person)
}

// @Summary Delete a person by id
// @Description Delete a person by id
// @Tags People
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Security ApiKeyAuth
// @Router /person/{id} [delete]
func DeletePerson(c *fiber.Ctx) error {
	id, err := GetValidId(c)
	if err != nil {
		return HandleInvalidId(c, err)
	}

	err = database.DAL.DeletePersonById(id)

	if err != nil {
		return ErorrAndDataResponse(c, err, id)
	}

	return DeletedResponse(c, "Person deleted", id)
}
