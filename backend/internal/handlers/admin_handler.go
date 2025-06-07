package handlers

import (
	"gochop/backend/internal/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LinkInfo represents the structure for link information
type LinkInfo struct {
	ID        int       `json:"id"`
	ShortCode string    `json:"short_code"`
	LongURL   string    `json:"long_url"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	ClickCount int      `json:"click_count"`
	UserID    string    `json:"user_id"`
}

// AnalyticsInfo represents analytics data for a specific link
type AnalyticsInfo struct {
	ShortCode       string                 `json:"short_code"`
	TotalClicks     int                    `json:"total_clicks"`
	ClicksByDate    []DailyClickData       `json:"clicks_by_date"`
	TopReferrers    []ReferrerData         `json:"top_referrers"`
	TopUserAgents   []UserAgentData        `json:"top_user_agents"`
	GeographicData  []GeographicData       `json:"geographic_data"`
}

// DailyClickData represents click data for a specific date
type DailyClickData struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

// ReferrerData represents referrer statistics
type ReferrerData struct {
	Referrer string `json:"referrer"`
	Clicks   int    `json:"clicks"`
}

// UserAgentData represents user agent statistics
type UserAgentData struct {
	UserAgent string `json:"user_agent"`
	Clicks    int    `json:"clicks"`
}

// GeographicData represents geographic statistics
type GeographicData struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Clicks  int    `json:"clicks"`
}

// GetAllLinks fetches all links with their click counts for the authenticated user
func GetAllLinks(c *fiber.Ctx) error {
	// Get user ID from context (set by NextAuth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	// Check if user is admin - if so, return all links, otherwise filter by user
	isAdmin, _ := c.Locals("isAdmin").(bool)
	
	var query string
	var args []interface{}
	
	if isAdmin {
		query = `
			SELECT l.id, l.short_code, l.long_url, l.context, l.created_at, l.expires_at, 
				   COALESCE(COUNT(a.id), 0) as click_count, l.user_id
			FROM links l
			LEFT JOIN analytics a ON l.short_code = a.short_code
			GROUP BY l.id, l.short_code, l.long_url, l.context, l.created_at, l.expires_at, l.user_id
			ORDER BY l.created_at DESC
		`
		args = []interface{}{}
	} else {
		query = `
			SELECT l.id, l.short_code, l.long_url, l.context, l.created_at, l.expires_at, 
				   COALESCE(COUNT(a.id), 0) as click_count, l.user_id
			FROM links l
			LEFT JOIN analytics a ON l.short_code = a.short_code
			WHERE l.user_id = $1
			GROUP BY l.id, l.short_code, l.long_url, l.context, l.created_at, l.expires_at, l.user_id
			ORDER BY l.created_at DESC
		`
		args = []interface{}{userID}
	}

	rows, err := db.DB.Query(db.Ctx, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch links",
		})
	}
	defer rows.Close()

	var links []LinkInfo
	for rows.Next() {
		var link LinkInfo
		err := rows.Scan(&link.ID, &link.ShortCode, &link.LongURL, &link.Context, &link.CreatedAt, &link.ExpiresAt, &link.ClickCount, &link.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not scan link data",
			})
		}
		links = append(links, link)
	}

	return c.JSON(links)
}

// GetAnalytics provides comprehensive analytics data for a specific link
func GetAnalytics(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	
	// Get user ID from context (set by NextAuth middleware)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User authentication required",
		})
	}

	// Check if user is admin
	isAdmin, _ := c.Locals("isAdmin").(bool)

	// Check if the link exists and user has access to it
	var exists bool
	var linkUserID string
	var checkQuery string
	var args []interface{}
	
	if isAdmin {
		// Admin can view any link
		checkQuery = "SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1), COALESCE((SELECT user_id FROM links WHERE short_code = $1), '')"
		args = []interface{}{shortCode, shortCode}
	} else {
		// Regular user can only view their own links
		checkQuery = "SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1 AND user_id = $2), COALESCE((SELECT user_id FROM links WHERE short_code = $1), '')"
		args = []interface{}{shortCode, userID, shortCode}
	}
	
	err := db.DB.QueryRow(db.Ctx, checkQuery, args...).Scan(&exists, &linkUserID)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Link not found or access denied",
		})
	}

	analytics := AnalyticsInfo{
		ShortCode: shortCode,
	}

	// Get total clicks
	totalClicksQuery := "SELECT COUNT(*) FROM analytics WHERE short_code = $1"
	err = db.DB.QueryRow(db.Ctx, totalClicksQuery, shortCode).Scan(&analytics.TotalClicks)
	if err != nil {
		analytics.TotalClicks = 0
	}

	// Get clicks by date (last 30 days)
	clicksByDateQuery := `
		SELECT DATE(clicked_at) as date, COUNT(*) as clicks
		FROM analytics 
		WHERE short_code = $1 AND clicked_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(clicked_at)
		ORDER BY date DESC
	`
	rows, err := db.DB.Query(db.Ctx, clicksByDateQuery, shortCode)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dailyData DailyClickData
			err := rows.Scan(&dailyData.Date, &dailyData.Clicks)
			if err == nil {
				analytics.ClicksByDate = append(analytics.ClicksByDate, dailyData)
			}
		}
	}

	// Get top referrers
	topReferrersQuery := `
		SELECT COALESCE(referrer, 'Direct') as referrer, COUNT(*) as clicks
		FROM analytics 
		WHERE short_code = $1
		GROUP BY referrer
		ORDER BY clicks DESC
		LIMIT 10
	`
	rows, err = db.DB.Query(db.Ctx, topReferrersQuery, shortCode)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var referrerData ReferrerData
			err := rows.Scan(&referrerData.Referrer, &referrerData.Clicks)
			if err == nil {
				analytics.TopReferrers = append(analytics.TopReferrers, referrerData)
			}
		}
	}

	// Get top user agents (simplified - just the first 50 chars for readability)
	topUserAgentsQuery := `
		SELECT LEFT(COALESCE(user_agent, 'Unknown'), 50) as user_agent, COUNT(*) as clicks
		FROM analytics 
		WHERE short_code = $1
		GROUP BY LEFT(user_agent, 50)
		ORDER BY clicks DESC
		LIMIT 10
	`
	rows, err = db.DB.Query(db.Ctx, topUserAgentsQuery, shortCode)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var userAgentData UserAgentData
			err := rows.Scan(&userAgentData.UserAgent, &userAgentData.Clicks)
			if err == nil {
				analytics.TopUserAgents = append(analytics.TopUserAgents, userAgentData)
			}
		}
	}

	// Get geographic data
	geographicQuery := `
		SELECT COALESCE(country, 'Unknown') as country, 
			   COALESCE(region, 'Unknown') as region,
			   COALESCE(city, 'Unknown') as city,
			   COUNT(*) as clicks
		FROM analytics 
		WHERE short_code = $1 AND country IS NOT NULL
		GROUP BY country, region, city
		ORDER BY clicks DESC
		LIMIT 20
	`
	rows, err = db.DB.Query(db.Ctx, geographicQuery, shortCode)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var geoData GeographicData
			err := rows.Scan(&geoData.Country, &geoData.Region, &geoData.City, &geoData.Clicks)
			if err == nil {
				analytics.GeographicData = append(analytics.GeographicData, geoData)
			}
		}
	}

	return c.JSON(analytics)
} 