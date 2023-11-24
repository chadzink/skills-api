package models

import (
	"time"
)

type PersonSkill struct {
	PersonID    uint      `json:"person_id" gorm:"primaryKey"`
	SkillID     uint      `json:"skill_id" gorm:"primaryKey"`
	ExpertiseID uint      `json:"expertise_id" gorm:"int;default:0"`
	LastUsed    time.Time `json:"last_used" gorm:"datetime;null;default:null"`

	Skill     Skill     `json:"skill" gorm:"foreignKey:SkillID" swaggerignore:"true"`
	Expertise Expertise `json:"expertise" gorm:"foreignKey:ExpertiseID" swaggerignore:"true"`
}

// Example of person skill in JSON
/*
[
	{
		"person_id": 1,
		"skill_id": 1,
		"expertise_id": 1,
		"last_used": "2021-01-01T00:00:00Z",
	}, {
		"person_id": 1,
		"skill_id": 2,
		"expertise_id": 2,
		"last_used": "2021-01-01T00:00:00Z",
	}, {
		"person_id": 1,
		"skill_id": 3,
		"expertise_id": 3,
		"last_used": "2021-01-01T00:00:00Z",
	}, {
		"person_id": 2,
		"skill_id": 1,
		"expertise_id": 1,
		"last_used": "2021-01-01T00:00:00Z",
	}, {
		"person_id": 2,
		"skill_id": 2,
		"expertise_id": 2,
		"last_used": "2021-01-01T00:00:00Z",
	}, {
		"person_id": 2,
		"skill_id": 3,
		"expertise_id": 3,
		"last_used": "2021-01-01T00:00:00Z",
	}
]
*/
