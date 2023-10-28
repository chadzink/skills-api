package database

import (
	"skills-api/models"
)

// Function to add associated skills to category
func (dal *DataAccessLayer) UpdatePersonSkillsForPerson(person *models.Person) error {
	if len(person.PersonSkills) > 0 {

		// Set the person id for each person skill
		for i, _ := range person.PersonSkills {
			person.PersonSkills[i].PersonID = person.ID
		}

		// Make a copy of the person skills to add so I can clear existing and add in the associated skills
		updatedPersonSkills := person.PersonSkills

		// remove existing associated skills
		if err := dal.Db.Model(&person).Association("PersonSkills").Clear(); err != nil {
			return err
		}

		dal.Db.Save(&person).Association("PersonSkills").Replace(updatedPersonSkills)

		// Load the relationships for each PersonSkill
		dal.LoadPersonSkillsForPerson(person)
	}

	return nil
}

// Load the Skills and Expertise for a Person
func (dal *DataAccessLayer) LoadPersonSkillsForPerson(person *models.Person) error {
	if len(person.PersonSkills) > 0 {
		for i, _ := range person.PersonSkills {
			// Load Skill
			dal.Db.Model(&person.PersonSkills[i]).Preload("Skill").First(&person.PersonSkills[i])

			// Load Expertise
			dal.Db.Model(&person.PersonSkills[i]).Preload("Expertise").First(&person.PersonSkills[i])
		}
	}

	return nil
}

// Create a new person entity and build the association with skills
func (dal *DataAccessLayer) CreatePerson(person *models.Person) error {
	err := dal.Db.Create(&person).Error
	if err != nil {
		return err
	}

	// Add associated person skills
	if err = dal.UpdatePersonSkillsForPerson(person); err != nil {
		return err
	}

	return nil
}

// Get all people and pre-load the skills
func (dal *DataAccessLayer) GetAllPeople() ([]models.Person, error) {
	var people []models.Person
	err := dal.Db.Model(&models.Person{}).Preload("PersonSkills").Find(&people).Error

	for i, _ := range people {
		dal.LoadPersonSkillsForPerson(&people[i])
	}

	return people, err
}

// Get a person by id and pre-load the skills
func (dal *DataAccessLayer) GetPersonById(id uint) (models.Person, error) {
	var person models.Person

	err := dal.Db.Model(&models.Person{}).Preload("PersonSkills").First(&person, id).Error
	dal.LoadPersonSkillsForPerson(&person)

	return person, err
}

// Update a person by id and rebuild the association with skills
func (dal *DataAccessLayer) UpdatePersonById(id uint, person *models.Person) error {
	var existingPerson models.Person
	err := dal.Db.Model(&models.Person{}).Preload("PersonSkills").First(&existingPerson, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Person{}).Where("id = ?", id).Updates(&person)

	// Set the id to the existing id
	person.ID = id

	// Add associated skills
	if err = dal.UpdatePersonSkillsForPerson(person); err != nil {
		return err
	}

	return nil
}

// Delete a person by id and remove the association with skills
func (dal *DataAccessLayer) DeletePersonById(id uint) error {
	var person models.Person
	err := dal.Db.Model(&models.Person{}).First(&person, id).Error
	if err != nil {
		return err
	}

	if err = dal.Db.Model(&person).Association("PersonSkills").Clear(); err != nil {
		return err
	}

	if err = dal.Db.Delete(&person).Error; err != nil {
		return err
	}

	return nil
}
