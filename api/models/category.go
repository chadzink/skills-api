package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"varchar(200);not null`
	Description string `json:"description" gorm:"text;null;default:null`
	ShortKey    string `json:"short_key" gorm:"varchar(10);not null;default:null`
	Active      bool   `json:"active" gorm:"bit;not null;default:1`

	Skills []*Skill `gorm:"many2many:skill_category;"`
}
