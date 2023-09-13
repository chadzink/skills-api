package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	Name        string `json:"name" gorm:"text;not null`
	Description string `json:"description" gorm:"text;null;default:null`
	ShortKey    string `json:"short_key" gorm:"text;not null;default:null`
	Active      bool   `json:"active" gorm:"bit;not null;default:1`

	Categories []*Category `gorm:"many2many:skill_category;"`

	// Additional fields that are not stored in the database
	CategoryIds []uint `json:"category_ids" gorm:"-"`
}

// Example of skills in JSON
/*
[
	{
		"name": "Go",
		"description": "Go is a compiled, statically typed programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
		"short_key": "go",
		"active": true,
	}, {
		"name": "JavaScript",
		"description": "JavaScript, often abbreviated as JS, is a programming language that conforms to the ECMAScript specification. JavaScript is high-level, often just-in-time compiled, and multi-paradigm. It has curly-bracket syntax, dynamic typing, prototype-based object-orientation, and first-class functions.",
		"short_key": "js",
		"active": true,
	}, {
		"name": "Python",
		"description": "Python is an interpreted, high-level and general-purpose programming language. Python's design philosophy emphasizes code readability with its notable use of significant indentation.",
		"short_key": "py",
		"active": true,
	}, {
		"name": "Java",
		"description": "Java is a class-based, object-oriented programming language that is designed to have as few implementation dependencies as possible. It is a general-purpose programming language intended to let application developers write once, run anywhere, meaning that compiled Java code can run on all platforms that support Java without the need for recompilation.",
		"short_key": "java",
		"active": true,
	}
]
*/
