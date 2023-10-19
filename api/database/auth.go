package database

import (
	"errors"

	"github.com/chadzink/skills-api/auth"
	"github.com/chadzink/skills-api/models"
)

const MAX_FAILED_ATTEMPTS = 5

// Find a user by email address and password
func (dal *DataAccessLayer) FindUserByCredentials(email, password string) (*models.User, error) {
	var user models.User

	// Find user by email
	if err := DAL.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	} else {
		// Check if the user is locked
		if user.Locked {
			return nil, errors.New("account locked")
		}

		// Check password
		if err := auth.CheckPassword(password, user.Password); err != nil {
			// We know the user matched, so we want to increase the failed attempts
			nextFailedAttempts := user.FailedAttempts + 1

			dal.Db.Model(&user).Update("failed_attempts", nextFailedAttempts)

			// If the user has failed MAX_FAILED_ATTEMPTS times, lock the account
			if nextFailedAttempts >= MAX_FAILED_ATTEMPTS {
				dal.Db.Model(&user).Update("locked", true)
			}

			return nil, errors.New("invalid password")
		}
	}

	return &user, nil
}

// Register a new user with email, disply name and password
func (dal *DataAccessLayer) RegisterNewUser(email, displayName, password string) (*models.User, error) {
	// Check if the email address is already in use
	var existingUser models.User
	if err := dal.Db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email address already in use")
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := models.User{
		Email:          email,
		DisplayName:    displayName,
		Password:       hashedPassword,
		FailedAttempts: 0,
		Locked:         false,
	}

	if err := dal.Db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
