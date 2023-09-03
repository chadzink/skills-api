package database

import (
	"github.com/chadzink/skills-api/models"
)

// Create a new skill entity and build the association with categories
func (dal *DataAccessLayer) CreateSkill(skill *models.Skill) error {
	err := dal.Db.Create(&skill).Error
	if err != nil {
		return err
	}

	for _, category := range skill.Categories {
		err = dal.Db.Model(&skill).Association("Categories").Append(&category)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get all skills and pre-load the categories
func (dal *DataAccessLayer) GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := dal.Db.Model(&models.Skill{}).Preload("Categories").Find(&skills).Error
	return skills, err
}

// Get a skill by id and pre-load the categories
func (dal *DataAccessLayer) GetSkillById(id uint) (models.Skill, error) {
	var skill models.Skill

	err := dal.Db.Model(&models.Skill{}).Preload("Categories").First(&skill, id).Error
	return skill, err
}

// Update a skill by id and rebuild the association with categories
func (dal *DataAccessLayer) UpdateSkillById(id uint, skill *models.Skill) error {
	var existingSkill models.Skill

	err := dal.Db.Model(&models.Skill{}).First(&existingSkill, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Skill{}).Where("id = ?", id).Updates(&skill)

	// Check if the category association has was passed in
	if len(skill.Categories) > 0 {
		// First remove the old associations
		err = dal.Db.Model(&skill).Association("Categories").Clear()
		if err != nil {
			return err
		}

		// Then add the new associations
		for _, category := range skill.Categories {
			err = dal.Db.Model(&skill).Association("Categories").Append(&category)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete a skill by id and remove the association with categories
func (dal *DataAccessLayer) DeleteSkillById(id uint) error {
	var skill models.Skill

	err := dal.Db.Model(&models.Skill{}).First(&skill, id).Error
	if err != nil {
		return err
	}

	err = dal.Db.Model(&skill).Association("Categories").Clear()
	if err != nil {
		return err
	}

	err = dal.Db.Delete(&skill).Error
	return err
}
