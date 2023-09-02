package database

import (
	"github.com/chadzink/skills-api/models"
)

// Get all skills and pre-load the categories
func (dbInstance *Dbinstance) GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := dbInstance.Db.Model(&models.Skill{}).Preload("Categories").Find(&skills).Error
	return skills, err
}

// Get a skill by id and pre-load the categories
func (dbInstance *Dbinstance) GetSkillById(id string) (models.Skill, error) {
	var skill models.Skill
	err := dbInstance.Db.Model(&models.Skill{}).Preload("Categories").First(&skill, id).Error
	return skill, err
}
