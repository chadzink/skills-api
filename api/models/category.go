package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"text;not null`
	Description string `json:"description" gorm:"text;null;default:null`
	ShortKey    string `json:"short_key" gorm:"varchar(10);not null;default:null`
	Active      bool   `json:"active" gorm:"bit;not null;default:1`

	Skills []*Skill `gorm:"many2many:skill_category;"`
}

// Example of categories in JSON
/*
[
	{
		"name": "Programming Languages",
		"description": "A programming language is a formal language comprising a set of instructions that produce various kinds of output. Programming languages are used in computer programming to implement algorithms.",
		"short_key": "pl",
		"active": true,
	}, {
		"name": "Databases",
		"description": "A database is an organized collection of structured information, or data, typically stored electronically in a computer system. A database is usually controlled by a database management system.",
		"short_key": "db",
		"active": true,
	}, {
		"name": "Operating Systems",
		"description": "An operating system is system software that manages computer hardware, software resources, and provides common services for computer programs.",
		"short_key": "os",
		"active": true,
	}, {
		"name": "Cloud Providers",
		"description": "A cloud provider is a company that offers some component of cloud computing -- typically infrastructure as a service, software as a service or platform as a service -- to other businesses or individuals.",
		"short_key": "cp",
		"active": true,
	}
]
*/
