package database

import (
	"skills-api/models"
)

// Create a new expertise entity
func (dal *DataAccessLayer) CreateExpertise(expertise *models.Expertise) error {
	err := dal.Db.Create(&expertise).Error
	if err != nil {
		return err
	}

	return nil
}

// Get all people and pre-load the skills
func (dal *DataAccessLayer) GetAllExpertise() ([]models.Expertise, error) {
	var expertises []models.Expertise
	err := dal.Db.Model(&models.Expertise{}).Find(&expertises).Error

	return expertises, err
}
