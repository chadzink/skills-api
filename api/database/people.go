package database

import (
	"github.com/chadzink/skills-api/models"
)

// Function to add associated skills to person
func (dal *DataAccessLayer) AddPersonSkillsForPerson(person *models.Person) error {

	// Check if the skill id list association has was passed in
	if len(person.PersonSkills) > 0 {
		// remove existing associated skills
		if err := dal.Db.Model(&person).Association("Skills").Clear(); err != nil {
			return err
		}

		// Add associated skills
		for i, personSkill := range person.PersonSkills {
			var skill models.Skill
			dal.Db.First(&skill, personSkill.SkillID)
			// Check if the skill was found by id
			if skill.ID > 0 {
				person.Skills = append(person.Skills, &skill)
				person.PersonSkills[i].PersonID = person.ID
			}
		}

		dal.Db.Save(&person).Association("Skills").Replace(person.Skills)
	}

	return nil
}

// Create a new person entity and build the association with skills
func (dal *DataAccessLayer) CreatePerson(person *models.Person) error {
	err := dal.Db.Create(&person).Error
	if err != nil {
		return err
	}

	// Add associated skills
	if err = dal.AddPersonSkillsForPerson(person); err != nil {
		return err
	}

	return nil
}

// Get all people and pre-load the skills
func (dal *DataAccessLayer) GetAllPeople() ([]models.Person, error) {
	var people []models.Person
	err := dal.Db.Model(&models.Person{}).Preload("Skills").Find(&people).Error

	// Build person skills lists
	for index, person := range people {
		for _, skill := range person.Skills {
			people[index].PersonSkills = append(people[index].PersonSkills, models.PersonSkill{
				PersonID:    person.ID,
				SkillID:     skill.ID,
				ExpertiseID: 0,
			})
		}
	}

	return people, err
}

// Get a person by id and pre-load the skills
func (dal *DataAccessLayer) GetPersonById(id uint) (models.Person, error) {
	var person models.Person

	err := dal.Db.Model(&models.Person{}).Preload("Skills").First(&person, id).Error

	// Build person skills list
	for _, skill := range person.Skills {
		person.PersonSkills = append(person.PersonSkills, models.PersonSkill{
			PersonID:    person.ID,
			SkillID:     skill.ID,
			ExpertiseID: 0,
		})
	}

	return person, err
}

// Update a person by id and rebuild the association with skills
func (dal *DataAccessLayer) UpdatePersonById(id uint, person *models.Person) error {
	var existingPerson models.Person
	err := dal.Db.Model(&models.Person{}).Preload("Skills").First(&existingPerson, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Person{}).Where("id = ?", id).Updates(&person)

	// Set the id to the existing id
	person.ID = id

	// Add associated skills
	if err = dal.AddPersonSkillsForPerson(person); err != nil {
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

	if err = dal.Db.Model(&person).Association("Skills").Clear(); err != nil {
		return err
	}

	if err = dal.Db.Delete(&person).Error; err != nil {
		return err
	}

	return nil
}
