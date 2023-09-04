package handlers

import (
	"github.com/chadzink/skills-api/database"
	"github.com/chadzink/skills-api/models"
	"github.com/gofiber/fiber/v2"
)

// CreatePerson creates a new person entity
func CreatePerson(c *fiber.Ctx) error {
	person := new(models.Person)

	if err := c.BodyParser(person); err != nil {
		return ErorrAndDataResponse(c, err, person)
	}

	database.DAL.CreatePerson(person)

	return DataResponse(c, person)
}

// ListPerson lists a person by id
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

// ListPeople lists all people
func ListPeople(c *fiber.Ctx) error {
	// TODO: Add pagination

	people, err := database.DAL.GetAllPeople()

	if err != nil {
		return ErorrAndDataResponse(c, err, people)
	}

	return DataResponse(c, people)
}

// UpdatePerson updates a person by id
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

// DeletePerson deletes a person by id
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
