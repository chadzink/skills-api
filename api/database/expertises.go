package database

import (
	"github.com/chadzink/skills-api/models"
	"gorm.io/gorm"
)

// Get all people and pre-load the skills
func (dal *DataAccessLayer) CreateDefaultExpertise() {
	dal.Db.FirstOrCreate(&models.Expertise{
		Model:       gorm.Model{ID: 1},
		Name:        "Beginner",
		Description: "A beginner is a person who is starting to learn or do something.",
		Order:       1,
	})
	dal.Db.FirstOrCreate(&models.Expertise{
		Model:       gorm.Model{ID: 2},
		Name:        "Intermediate",
		Description: "An intermediate is a person who has a level of knowledge or skill between a beginner and an expert.",
		Order:       2,
	})
	dal.Db.FirstOrCreate(&models.Expertise{
		Model:       gorm.Model{ID: 3},
		Name:        "Advanced",
		Description: "An advanced is a person who is very skilled or highly trained in a particular field.",
		Order:       3,
	})
	dal.Db.FirstOrCreate(&models.Expertise{
		Model:       gorm.Model{ID: 4},
		Name:        "Expert",
		Description: "An expert is a person who is very knowledgeable about or skilful in a particular area.",
		Order:       4,
	})
	dal.Db.FirstOrCreate(&models.Expertise{
		Model:       gorm.Model{ID: 5},
		Name:        "N/A",
		Description: "N/A is used when the level of expertise is not applicable.",
		Order:       5,
	})
}

// Get all people and pre-load the skills
func (dal *DataAccessLayer) GetAllExpertise() ([]models.Expertise, error) {
	var expertises []models.Expertise
	err := dal.Db.Model(&models.Expertise{}).Find(&expertises).Error

	return expertises, err
}
