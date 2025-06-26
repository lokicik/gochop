package handlers

import (
	"context"
	"gochop/backend/internal/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck responds to health check requests
func HealthCheck(c *fiber.Ctx) error {
	// Check database connection
	dbStatus := "healthy"
	if db.DB == nil {
		dbStatus = "unhealthy"
	} else {
		// Try to ping the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := db.DB.Ping(ctx); err != nil {
			dbStatus = "unhealthy"
		}
	}

	// Check Redis connection
	redisStatus := "healthy"
	if db.RedisClient == nil {
		redisStatus = "unhealthy"
	} else {
		// Try to ping Redis
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := db.RedisClient.Ping(ctx).Err(); err != nil {
			redisStatus = "unhealthy"
		}
	}

	status := "healthy"
	statusCode := fiber.StatusOK
	if dbStatus == "unhealthy" || redisStatus == "unhealthy" {
		status = "unhealthy"
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"database":  dbStatus,
		"redis":     redisStatus,
	})
} 