package models

import (
	"gorm.io/gorm"
)

// Type login request with username and password (username is email address)
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Type login response with JWT token
type LoginResponse struct {
	Token string `json:"token"`
}

// Type for register request
type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

// Type for a user
type User struct {
	gorm.Model     `json:"-" swaggerignore:"true"`
	DisplayName    string `json:"display_name" gorm:"text;not null"`
	Email          string `json:"email" gorm:"text;not null;unique"`
	Password       string `json:"password" gorm:"text;not null"`
	FailedAttempts int    `json:"failed_attempts" gorm:"int;not null;default:0"`
	Locked         bool   `json:"locked" gorm:"bool;not null;default:false"`
}
