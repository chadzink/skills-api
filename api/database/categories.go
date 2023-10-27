package database

import (
	"skills-api/models"
)

// Function to add associated skills to category
func (dal *DataAccessLayer) UpdateSkillsForCategory(category *models.Category) error {
	// Check if the skill id list association has was passed in
	if len(category.SkillIds) > 0 {
		for _, skillId := range category.SkillIds {
			var skill models.Skill
			dal.Db.First(&skill, skillId)

			if skill.ID > 0 {
				category.Skills = append(category.Skills, &skill)
			}
		}
	}

	if len(category.Skills) > 0 {
		// Mkae a copy of the skills to add so I can clear existing and add in the associated skills
		updatedCategorySkills := category.Skills

		// remove existing associated skills
		if err := dal.Db.Model(&category).Association("Skills").Clear(); err != nil {
			return err
		}

		dal.Db.Save(&category).Association("Skills").Replace(updatedCategorySkills)
	}

	return nil
}

// Create a new category entity and build the association with skills
func (dal *DataAccessLayer) CreateCategory(category *models.Category) error {
	err := dal.Db.Create(&category).Error
	if err != nil {
		return err
	}

	if err = dal.UpdateSkillsForCategory(category); err != nil {
		return err
	}

	return nil
}

// Get all categories and pre-load the skills
func (dal *DataAccessLayer) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := dal.Db.Model(&models.Category{}).Preload("Skills").Find(&categories).Error

	// Build skill id lists
	for index, category := range categories {
		for _, skill := range category.Skills {
			categories[index].SkillIds = append(categories[index].SkillIds, skill.ID)
		}
	}

	return categories, err
}

// Get a category by id and pre-load the skills
func (dal *DataAccessLayer) GetCategoryById(id uint) (models.Category, error) {
	var category models.Category

	err := dal.Db.Model(&models.Category{}).Preload("Skills").First(&category, id).Error
	// Build skill id list
	for _, skill := range category.Skills {
		category.SkillIds = append(category.SkillIds, skill.ID)
	}

	return category, err
}

// Update a category by id and rebuild the association with skills
func (dal *DataAccessLayer) UpdateCategoryById(id uint, category *models.Category) error {
	var existingCategory models.Category
	err := dal.Db.Model(&models.Category{}).First(&existingCategory, id).Error
	if err != nil {
		return err
	}

	dal.Db.Model(&models.Category{}).Where("id = ?", id).Updates(&category)

	// Set the id to the existing id
	category.ID = id

	// Check if the skill association has was passed in
	if err = dal.UpdateSkillsForCategory(category); err != nil {
		return err
	}

	return nil
}

// Delete a category by id and remove the association with skills
func (dal *DataAccessLayer) DeleteCategoryById(id uint) error {
	var category models.Category
	err := dal.Db.Model(&models.Category{}).First(&category, id).Error
	if err != nil {
		return err
	}

	if err = dal.Db.Model(&category).Association("Skills").Clear(); err != nil {
		return err
	}

	if err = dal.Db.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
