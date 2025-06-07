package handlers

import (
	"gochop/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest represents the login request structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response structure
type LoginResponse struct {
	Token   string `json:"token"`
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Message string `json:"message"`
}

// Login handles user authentication and JWT token generation
// Note: This is a simplified demo implementation
// In production, you'd want proper password hashing, user database, etc.
func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Simple demo authentication (replace with real user verification)
	var userID string
	var isAdmin bool

	switch {
	case req.Username == "admin" && req.Password == "admin123":
		userID = "admin-user"
		isAdmin = true
	case req.Username == "user" && req.Password == "user123":
		userID = "regular-user"
		isAdmin = false
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(userID, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(LoginResponse{
		Token:   token,
		UserID:  userID,
		IsAdmin: isAdmin,
		Message: "Login successful",
	})
}

// GetProfile returns the current user's profile information
func GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	isAdmin, ok := c.Locals("isAdmin").(bool)
	if !ok {
		isAdmin = false // Default to false if not set
	}

	return c.JSON(fiber.Map{
		"user_id":  userID,
		"is_admin": isAdmin,
		"message":  "Profile retrieved successfully",
	})
}

// GenerateAdminToken generates an admin token (for development/testing purposes)
func GenerateAdminToken(c *fiber.Ctx) error {
	token, err := middleware.GenerateToken("dev-admin", true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token":   token,
		"message": "Admin token generated for development",
		"note":    "This endpoint should be removed in production",
	})
} 