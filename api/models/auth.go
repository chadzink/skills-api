package models

import (
	"time"

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

// Type for the new API key request
type NewAPIKeyRequest struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	ExpiresOn time.Time `json:"expires_on"`
}

// Type for the new API key request
type VerifyAPIKeyRequest struct {
	Email string `json:"email"`
	Key   string `json:"key"`
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

type UserAPIKey struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	UserID     uint      `json:"user_id" gorm:"int;not null"`
	Key        string    `json:"key" gorm:"text;not null;unique"`
	ExpiresOn  time.Time `json:"expires_on" gorm:"timestamp;not null"`
}
