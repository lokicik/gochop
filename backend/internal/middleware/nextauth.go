package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// NextAuthClaims represents the claims structure from NextAuth.js JWT
type NextAuthClaims struct {
	Sub       string `json:"sub"`        // User ID
	Name      string `json:"name"`       // User name
	Email     string `json:"email"`      // User email
	Picture   string `json:"picture"`    // User picture URL
	IsAdmin   bool   `json:"isAdmin"`    // Admin flag
	Iat       int64  `json:"iat"`        // Issued at
	Exp       int64  `json:"exp"`        // Expiration
	Jti       string `json:"jti"`        // JWT ID
	jwt.RegisteredClaims
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS represents a JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// NextAuthMiddleware validates NextAuth.js JWT tokens
func NextAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check for Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Parse and validate the token
		userID, isAdmin, err := validateNextAuthToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: " + err.Error(),
			})
		}

		// Store user information in context
		c.Locals("userID", userID)
		c.Locals("isAdmin", isAdmin)

		return c.Next()
	}
}

// AdminOnlyMiddleware ensures only admin users can access the endpoint
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("isAdmin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}

// validateNextAuthToken validates a NextAuth.js JWT token
func validateNextAuthToken(tokenString string) (userID string, isAdmin bool, err error) {
	// For development, we'll use a simple validation
	// In production, you should validate against NextAuth.js JWKS endpoint
	
	secret := os.Getenv("NEXTAUTH_SECRET")
	if secret == "" {
		return "", false, fmt.Errorf("NEXTAUTH_SECRET not configured")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &NextAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", false, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*NextAuthClaims)
	if !ok || !token.Valid {
		return "", false, fmt.Errorf("invalid token claims")
	}

	// Check expiration
	if claims.Exp < time.Now().Unix() {
		return "", false, fmt.Errorf("token expired")
	}

	return claims.Sub, claims.IsAdmin, nil
}

// validateJWKS validates token using JWKS (for production use)
func validateJWKS(tokenString string) (userID string, isAdmin bool, err error) {
	// This is a more secure approach for production
	// You would fetch the JWKS from NextAuth.js endpoint
	// and validate the token signature against it
	
	nextAuthURL := os.Getenv("NEXTAUTH_URL")
	if nextAuthURL == "" {
		return "", false, fmt.Errorf("NEXTAUTH_URL not configured")
	}

	// Parse token to get key ID
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &NextAuthClaims{})
	if err != nil {
		return "", false, fmt.Errorf("failed to parse token: %w", err)
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return "", false, fmt.Errorf("key ID not found in token header")
	}

	// Fetch JWKS
	jwksURL := fmt.Sprintf("%s/api/auth/jwks", nextAuthURL)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return "", false, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Find the key
	var publicKey *rsa.PublicKey
	for _, key := range jwks.Keys {
		if key.Kid == keyID {
			publicKey, err = jwkToRSAPublicKey(key)
			if err != nil {
				return "", false, fmt.Errorf("failed to convert JWK to RSA key: %w", err)
			}
			break
		}
	}

	if publicKey == nil {
		return "", false, fmt.Errorf("public key not found for key ID: %s", keyID)
	}

	// Validate token with public key
	validatedToken, err := jwt.ParseWithClaims(tokenString, &NextAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return "", false, fmt.Errorf("failed to validate token: %w", err)
	}

	claims, ok := validatedToken.Claims.(*NextAuthClaims)
	if !ok || !validatedToken.Valid {
		return "", false, fmt.Errorf("invalid token claims")
	}

	return claims.Sub, claims.IsAdmin, nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert to big integers
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Create RSA public key
	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
} 