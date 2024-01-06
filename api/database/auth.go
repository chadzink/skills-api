package database

import (
	"errors"
	"time"

	"skills-api/auth"
	"skills-api/models"
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

// The CreateUserAPIKey method generates a unique API key string for the user based on the user ID
func (dal *DataAccessLayer) CreateUserAPIKey(requestBy *models.NewAPIKeyRequest) (*models.UserAPIKey, error) {
	// Generate a new API key
	apiKey := auth.GenerateRandomUTF8String(50)

	// Find the user by email address and password
	user, err := dal.FindUserByCredentials(requestBy.Email, requestBy.Password)
	if err != nil {
		return nil, err
	}

	// Check if the API key already exists
	existingUserApiKey := &models.UserAPIKey{
		UserID:    user.ID,
		ExpiresOn: requestBy.ExpiresOn,
		Key:       apiKey,
	}
	if err := dal.Db.Where(existingUserApiKey).First(&existingUserApiKey).Error; err == nil {
		return nil, errors.New("api key already exists")
	}

	// Add the API key to the database for the user
	if err := dal.Db.Create(existingUserApiKey).Error; err != nil {
		return nil, err
	}

	return existingUserApiKey, nil
}

// Check a UserAPIKey method takes a string values and checks if it is a valid API key, then returns the users
func (dal *DataAccessLayer) VerifyUserAPIKey(verifyRequest models.VerifyAPIKeyRequest) (*models.UserAPIKey, error) {
	// Find user by email
	var user models.User
	if err := dal.Db.Where("email = ?", verifyRequest.Email).First(&user).Error; err != nil {
		return nil, err
	}

	existingUserApiKey := &models.UserAPIKey{
		UserID: user.ID,
		Key:    verifyRequest.Key,
	}

	// Find users API key
	if err := dal.Db.Where(existingUserApiKey).First(&existingUserApiKey).Error; err != nil {
		return nil, err
	}

	// Check if the API key is expired
	if existingUserApiKey.ExpiresOn.Before(time.Now()) {
		return nil, errors.New("api key expired")
	}

	return existingUserApiKey, nil
}

// List all API keys for a user
func (dal *DataAccessLayer) ListUserAPIKeys(userID uint) ([]models.UserAPIKey, error) {
	var userAPIKeys []models.UserAPIKey

	if err := dal.Db.Where("user_id = ?", userID).Find(&userAPIKeys).Error; err != nil {
		return nil, err
	}

	return userAPIKeys, nil
}

// Get a user by Email
func (dal *DataAccessLayer) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := dal.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
