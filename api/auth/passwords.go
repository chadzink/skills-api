package auth

import (
	"crypto/md5"
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
