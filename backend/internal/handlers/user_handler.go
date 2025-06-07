package handlers

import (
	"gochop/backend/internal/db"
	"gochop/backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

var userService = services.NewUserService()

// GetUserProfile retrieves the current user's profile information with full details
func GetUserProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by NextAuth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	// Get user from database
	user, err := userService.GetUser(db.Ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user statistics
	stats, err := userService.GetUserStats(db.Ctx, userID)
	if err != nil {
		// If stats fail, continue without them
		stats = map[string]interface{}{}
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"stats": stats,
	})
}

// UpdateProfile updates the current user's profile information
func UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by NextAuth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	// Parse request body
	var input services.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update user
	user, err := userService.UpdateUser(db.Ctx, userID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// GetUserStats retrieves statistics for the current user
func GetUserStats(c *fiber.Ctx) error {
	// Get user ID from context (set by NextAuth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	// Get user statistics
	stats, err := userService.GetUserStats(db.Ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve statistics",
		})
	}

	return c.JSON(stats)
}

// CreateOrUpdateUser creates or updates a user (called by NextAuth webhook or during authentication)
func CreateOrUpdateUser(c *fiber.Ctx) error {
	// Parse request body
	var input services.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if input.ID == "" || input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID and email are required",
		})
	}

	// Create or update user
	user, err := userService.CreateUser(db.Ctx, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create or update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User created/updated successfully",
		"user":    user,
	})
}

// ListUsers returns all users (admin only)
func ListUsers(c *fiber.Ctx) error {
	// Get all users
	users, err := userService.ListUsers(db.Ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(users)
}

// GetUserByID retrieves a specific user by ID (admin only)
func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get user from database
	user, err := userService.GetUser(db.Ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user statistics
	stats, err := userService.GetUserStats(db.Ctx, userID)
	if err != nil {
		stats = map[string]interface{}{}
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"stats": stats,
	})
} 