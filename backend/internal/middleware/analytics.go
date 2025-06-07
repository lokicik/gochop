package middleware

import (
	"gochop/backend/internal/db"
	"gochop/backend/internal/services"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsData represents the data to be logged for analytics
type AnalyticsData struct {
	ShortCode string
	IPAddress string
	UserAgent string
	Referrer  string
	Country   string
	Region    string
	City      string
}

// LogAnalytics logs analytics data asynchronously to avoid blocking the request
func LogAnalytics(data AnalyticsData) {
	go func() {
		insertSQL := `INSERT INTO analytics (short_code, ip_address, user_agent, referrer, country, region, city) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := db.DB.Exec(db.Ctx, insertSQL, data.ShortCode, data.IPAddress, data.UserAgent, data.Referrer, data.Country, data.Region, data.City)
		if err != nil {
			// Log error but don't fail the request
			// In a production environment, you'd want proper logging here
		}
	}()
}

// GetClientIP extracts the real client IP from the request
func GetClientIP(c *fiber.Ctx) string {
	// Check for X-Forwarded-For header (proxy/load balancer)
	xForwardedFor := c.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// Check for X-Real-IP header
	xRealIP := c.Get("X-Real-IP")
	if xRealIP != "" && net.ParseIP(xRealIP) != nil {
		return xRealIP
	}

	// Fallback to connection remote address
	return c.IP()
}

// AnalyticsMiddleware middleware for logging analytics data
func AnalyticsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only log analytics for shortlink redirects (not API endpoints)
		path := c.Path()
		if strings.HasPrefix(path, "/api/") {
			return c.Next()
		}

		// Extract shortCode from the path (assuming it's the first segment)
		shortCode := strings.TrimPrefix(path, "/")
		if shortCode == "" || shortCode == "favicon.ico" {
			return c.Next()
		}

		// Get geographic data from IP
		clientIP := GetClientIP(c)
		geoData, err := services.GetLocationFromIP(clientIP)
		if err != nil {
			// Use fallback if geo service fails
			geoData = services.GetLocationFromIPFallback(clientIP)
		}

		// Log analytics data
		LogAnalytics(AnalyticsData{
			ShortCode: shortCode,
			IPAddress: clientIP,
			UserAgent: c.Get("User-Agent"),
			Referrer:  c.Get("Referer"),
			Country:   geoData.Country,
			Region:    geoData.Region,
			City:      geoData.City,
		})

		return c.Next()
	}
} 