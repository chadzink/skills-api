package models

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name    string `json:"name" gorm:"text;not null`
	Email   string `json:"email" gorm:"text;null;default:null`
	Phone   string `json:"phone" gorm:"text;null;default:null`
	Profile string `json:"profile" gorm:"text;null;default:null`

	Skills []*Skill `gorm:"many2many:person_skills;"`

	// Additional fields that are not stored in the database
	PersonSkills []PersonSkill `json:"person_skills" gorm:"-"`
}

type PersonSkill struct {
	PersonID    uint      `json:"person_id" gorm:"primaryKey"`
	SkillID     uint      `json:"skill_id" gorm:"primaryKey"`
	ExpertiseID uint      `json:"expertise_id" gorm:"primaryKey"`
	LastUsed    time.Time `json:"last_used" gorm:"primaryKey"`
}

// Example of categories in JSON
/*
[
	{
		"name": "John",
		"email": "john@email.com",
		"phone": "555-555-5555",
		"profile": "John is a software developer with 10 years of experience.",
	}, {
		"name": "Jane",
		"email": "jane@email.com",
		"phone": "555-555-5555",
		"profile": "Jane is a software developer with 15 years of experience.",
	}, {
		"name": "Joe",
		"email": "joe@email.com",
		"phone": "555-555-5555",
		"profile": "Joe is a software developer with 20 years of experience.",
	},
]
*/
