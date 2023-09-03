package database

import (
	"github.com/chadzink/skills-api/models"
)

// Create a new category entity and build the association with skills
func (dal *DataAccessLayer) CreateCategory(category *models.Category) error {
	err := dal.Db.Create(&category).Error
	if err != nil {
		return err
	}

	for _, skill := range category.Skills {
		err = dal.Db.Model(&category).Association("Skills").Append(&skill)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get all categories and pre-load the skills
func (dal *DataAccessLayer) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := dal.Db.Model(&models.Category{}).Preload("Skills").Find(&categories).Error
	return categories, err
}

// Get a category by id and pre-load the skills
func (dal *DataAccessLayer) GetCategoryById(id uint) (models.Category, error) {
	var category models.Category

	err := dal.Db.Model(&models.Category{}).Preload("Skills").First(&category, id).Error
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

	// Check if the skill association has was passed in
	if len(category.Skills) > 0 {
		dal.Db.Save(&category).Association("Skills").Replace(category.Skills)
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
