package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

//TO DO: Make these stronger with bcrypt and add a salt

// HashPassword hashes a plain text password using MD5
func HashPassword(password string) (string, error) {
	// Generate an MD5 hash of the plain text password
	hash := md5.Sum([]byte(password))

	// Encode the hash as a hexadecimal string
	encoded := hex.EncodeToString(hash[:])

	return encoded, nil
}

// CheckPassword checks a plain text password against an encoded password
func CheckPassword(password, encoded string) error {
	// Generate an MD5 hash of the plain text password
	hash := md5.Sum([]byte(password))

	// Encode the hash as a hexadecimal string
	encodedHash := hex.EncodeToString(hash[:])

	// Compare the generated hash to the stored hash
	if encodedHash != encoded {
		return errors.New("invalid password")
	}

	return nil
}

// The GenerateRandomString method generates a random string of a specified length
func GenerateRandomString(length int) (string, error) {

	generateRandomBytes := func(n int) ([]byte, error) {
		// Create a byte slice of the specified length
		b := make([]byte, n)

		// Generate random bytes and store them in the byte slice
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}

		return b, nil

	}

	// Generate a random string of a specified length
	randomString, err := generateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return string(randomString), nil
}
