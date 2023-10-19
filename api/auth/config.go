package auth

import "os"

// The secret key used to sign the JWT
var JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
