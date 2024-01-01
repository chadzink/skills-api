package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	mathRand "math/rand"
	"time"
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

func GenerateRandomUTF8String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"1234567890"

	seededRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
