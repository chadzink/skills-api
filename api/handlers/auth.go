package handlers

import (
	"time"

	"skills-api/auth"

	"skills-api/database"
	"skills-api/models"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
)

func CreateJWTClaims(user *models.User) jtoken.MapClaims {
	day := time.Hour * 24

	return jtoken.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(day * 24).Unix(),
	}
}

func CreateJWTToken(user *models.User) (string, error) {
	// Create the JWT claims, which includes the user ID and expiry time
	claims := CreateJWTClaims(user)

	// Create token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(auth.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

// @Summary User Login
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login JSON object that needs to be created"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	// Extract the credentials from the request body
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Find the user by credentials
	user, err := database.DAL.FindUserByCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Create token
	token, err := CreateJWTToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Return the token
	return c.JSON(models.LoginResponse{
		Token: token,
	})
}

// @Summary New User Register
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body models.RegisterRequest true "Register JSON object that needs to be created"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func RegisterNewUser(c *fiber.Ctx) error {
	// Extract the credentials from the request body
	registerRequest := new(models.RegisterRequest)
	if err := c.BodyParser(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Register the user
	user, err := database.DAL.RegisterNewUser(registerRequest.Email, registerRequest.DisplayName, registerRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Create token
	token, err := CreateJWTToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Return the token
	return c.JSON(models.LoginResponse{
		Token: token,
	})
}
