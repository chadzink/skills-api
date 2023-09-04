package database

import (
	"github.com/chadzink/skills-api/models"
)

// Create a new person entity and build the association with skills
func (dal *DataAccessLayer) CreatePerson(person *models.Person) error {
	err := dal.Db.Create(&person).Error
	if err != nil {
		return err
	}

	return nil
}

// Get all people and pre-load the skills
func (dal *DataAccessLayer) GetAllPeople() ([]models.Person, error) {
	var people []models.Person
	err := dal.Db.Model(&models.Person{}).Preload("Skills").Find(&people).Error

	return people, err
}

// Get a person by id and pre-load the skills
func (dal *DataAccessLayer) GetPersonById(id uint) (models.Person, error) {
	var person models.Person

	err := dal.Db.Model(&models.Person{}).Preload("Skills").First(&person, id).Error
	return person, err
}

// Update a person by id and rebuild the association with skills
func (dal *DataAccessLayer) UpdatePersonById(id uint, person *models.Person) error {
	var existingPerson models.Person
	err := dal.Db.Model(&models.Person{}).First(&existingPerson, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Person{}).Where("id = ?", id).Updates(&person)

	return nil
}

// Delete a person by id and remove the association with skills
func (dal *DataAccessLayer) DeletePersonById(id uint) error {
	var person models.Person
	err := dal.Db.Model(&models.Person{}).First(&person, id).Error
	if err != nil {
		return err
	}

	if err = dal.Db.Model(&person).Association("Skills").Clear(); err != nil {
		return err
	}

	if err = dal.Db.Delete(&person).Error; err != nil {
		return err
	}

	return nil
}
