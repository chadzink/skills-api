package handlers

import (
	"fmt"
	"strings"
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

func hasJWTToken(c *fiber.Ctx) bool {
	return c.Get("Authorization") != ""
}

func hasAPIKey(c *fiber.Ctx) bool {
	return c.Get("X-API-Email") != "" && c.Get("X-API-Key") != ""
}

func ExtractTokenFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")

	// Check if the user has an bearer token in the Authorization header
	if bearerToken == "" {
		return "", fmt.Errorf("Authorization header is required")
	}

	// Check if the token is in the correct format
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}
	// Extract the actual token
	return strings.TrimPrefix(bearerToken, "Bearer "), nil
}

func ParseToken(tokenString string) (*jtoken.Token, error) {
	token, err := jtoken.Parse(tokenString, func(token *jtoken.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jtoken.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.JWTSecretKey), nil
	})
	return token, err
}

func ExtractClaims(tokenString string) (jtoken.MapClaims, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jtoken.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
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

// Create UserAPIKey handler, this generated a unique API key string for the user based on the user ID
// @Summary Create a new API Key
// @Description Create a new API Key
// @Tags Auth
// @Accept json
// @Produce json
// @Param apiKeyRequest body models.NewAPIKeyRequest true "Create API Key JSON object that needs to be created"
// @Success 200 {object} ResponseResult[[]models.UserAPIKey]
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/api_key [post]
func CreateAPIKey(c *fiber.Ctx) error {

	apiKeyRequest := new(models.NewAPIKeyRequest)
	if err := c.BodyParser(apiKeyRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Create the API Key
	apiKey, err := database.DAL.CreateUserAPIKey(apiKeyRequest)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return DataResponse(c, apiKey)
}

// Get the current user verified by the login process for JWT or API Key
// @Summary Get the current user
// @Description Get the current user
// @Tags Auth
// @Produce json
// @Success 200 {object} ResponseResult[models.User]
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /user [get]
func GetCurrentUser(c *fiber.Ctx) error {
	// Check if the user has an bearer token in the Authorization header
	var userEmail string

	if hasJWTToken(c) {
		// Extract the user from the Authorization header
		tokenString, err := ExtractTokenFromHeader(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Extract the claims from the token
		claims, err := ExtractClaims(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Get the user from the token claims
		userEmail = claims["email"].(string)
	} else if hasAPIKey(c) {
		// Get the user from the API key
		userEmail = c.Get("X-API-Email")

		// Check if the API key and email are valid
		if _, err := database.DAL.VerifyUserAPIKey(models.VerifyAPIKeyRequest{
			Email: userEmail,
			Key:   c.Get("X-API-Key"),
		}); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// If we get this fgar we have a valid user email
	user, err := database.DAL.FindUserByEmail(userEmail)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return DataResponse(c, user)
}

// List all API keys for the authenticated user
// @Summary List all API keys for the authenticated user
// @Description List all API keys for the current user
// @Tags Auth
// @Produce json
// @Success 200 {object} ResponseResult[[]models.UserAPIKey]
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /user/api_key [get]
// func ListAPIKeys(c *fiber.Ctx) error {
// 	// Get the user from the JWT token
// 	user := c.Locals("user").(*models.User)

// 	// Get the API keys for the user
// 	apiKeys, err := database.DAL.GetUserAPIKeys(user)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return DataResponse(c, apiKeys)
// }
