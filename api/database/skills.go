package database

import (
	"skills-api/models"
)

// Function to add associated categories to skill
func (dal *DataAccessLayer) UpdateCategoriesForSkill(skill *models.Skill) error {
	// Check if the category id list association has was passed in
	if len(skill.CategoryIds) > 0 {
		// Add passed in associated categories
		for _, categoryId := range skill.CategoryIds {
			var category models.Category
			dal.Db.First(&category, categoryId)
			if category.ID > 0 {
				skill.Categories = append(skill.Categories, &category)
			}
		}
	}

	if len(skill.Categories) > 0 {
		// Mkae a copy of the categories to add so I can clear existing and add in the associated categories
		updatedSkillCategories := skill.Categories

		// remove existing associated skills
		if err := dal.Db.Model(&skill).Association("Categories").Clear(); err != nil {
			return err
		}

		dal.Db.Save(&skill).Association("Categories").Replace(updatedSkillCategories)
	}

	return nil
}

// Create a new skill entity and build the association with categories
func (dal *DataAccessLayer) CreateSkill(skill *models.Skill) error {
	err := dal.Db.Create(&skill).Error
	if err != nil {
		return err
	}

	if err = dal.UpdateCategoriesForSkill(skill); err != nil {
		return err
	}

	return nil
}

// Get all skills and pre-load the categories
func (dal *DataAccessLayer) GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := dal.Db.Model(&models.Skill{}).Preload("Categories").Find(&skills).Error

	// Build category id lists
	for index, skill := range skills {
		for _, category := range skill.Categories {
			skills[index].CategoryIds = append(skills[index].CategoryIds, category.ID)
		}
	}

	return skills, err
}

// Get a skill by id and pre-load the categories
func (dal *DataAccessLayer) GetSkillById(id uint) (models.Skill, error) {
	var skill models.Skill

	err := dal.Db.Model(&models.Skill{}).Preload("Categories").First(&skill, id).Error

	// Build category id lists
	for _, category := range skill.Categories {
		skill.CategoryIds = append(skill.CategoryIds, category.ID)
	}

	return skill, err
}

// Update a skill by id and rebuild the association with categories
func (dal *DataAccessLayer) UpdateSkillById(id uint, skill *models.Skill) error {
	// Make sure the skill exist in the database
	var existingSkill models.Skill
	err := dal.Db.Model(&models.Skill{}).First(&existingSkill, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Skill{}).Where("id = ?", id).Updates(&skill)

	// Set the id to the existing id
	skill.ID = id

	if err = dal.UpdateCategoriesForSkill(skill); err != nil {
		return err
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
