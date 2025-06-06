package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	// Get JWT secret from environment variable or generate a random one
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Generate a random secret for development (should be set in production)
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		secret = hex.EncodeToString(randomBytes)
	}
	jwtSecret = []byte(secret)
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID string, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours expiration
	claims := &Claims{
		UserID: userID,
		Admin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gochop",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// JWTMiddleware middleware to protect routes with JWT authentication
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store claims in context for use in handlers
		c.Locals("userID", claims.UserID)
		c.Locals("isAdmin", claims.Admin)

		return c.Next()
	}
}

// AdminMiddleware middleware to ensure the user has admin privileges
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("isAdmin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin privileges required",
			})
		}
		return c.Next()
	}
}

// OptionalJWTMiddleware middleware that validates JWT if present but doesn't require it
func OptionalJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ValidateToken(tokenString)
			if err == nil {
				c.Locals("userID", claims.UserID)
				c.Locals("isAdmin", claims.Admin)
			}
		}
		return c.Next()
	}
} 