package models

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	Name       string `json:"name" gorm:"text;not null`
	Email      string `json:"email" gorm:"text;null;default:null`
	Phone      string `json:"phone" gorm:"text;null;default:null`
	Profile    string `json:"profile" gorm:"text;null;default:null`

	PersonSkills []*PersonSkill `json:"person_skills" gorm:"foreignKey:PersonID" swaggerignore:"true"`
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
