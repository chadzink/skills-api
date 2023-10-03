package models

import (
	"gorm.io/gorm"
)

type Expertise struct {
	gorm.Model
	Name        string `json:"name" gorm:"text;not null`
	Description string `json:"description" gorm:"text;null;default:null`
	Order       int    `json:"order" gorm:"int;not null;default:0`
}

// Example of categories in JSON
/*
[
	{
		"name": "Beginner",
		"description": "A beginner is a person who is starting to learn or do something.",
		"order": "1",
	}, {
		"name": "Intermediate",
		"description": "An intermediate is a person who has a level of knowledge or skill between a beginner and an expert.",
		"order": "2",
	}, {
		"name": "Advanced",
		"description": "An advanced is a person who is very skilled or highly trained in a particular field.",
		"order": "3",
	}, {
		"name": "Expert",
		"description": "An expert is a person who is very knowledgeable about or skilful in a particular area.",
		"order": "4",
	}, {
		"name": "N/A",
		"description": "N/A is used when the level of expertise is not applicable.",
		"order": "5",
	}
]
*/
