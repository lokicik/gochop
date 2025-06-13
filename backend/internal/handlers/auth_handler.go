package handlers

import (
	"github.com/gofiber/fiber/v2"
)

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