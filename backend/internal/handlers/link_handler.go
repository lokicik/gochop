package handlers

import (
	"crypto/rand"
	"fmt"
	"gochop/backend/internal/db"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

const (
	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeLength   = 6
	cacheDuration     = 6 * time.Hour
	defaultExpiration = 90 * 24 * time.Hour // 90 days
)

// getBaseURL returns the base URL for short links from environment or default
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001" // Default for development
	}
	return baseURL
}

// validateURL checks if the provided URL is valid
func validateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include a scheme (http:// or https://)")
	}
	
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a host")
	}
	
	return nil
}

// validateAlias checks if the provided alias is valid
func validateAlias(alias string) error {
	if alias == "" {
		return nil // Empty alias is allowed
	}
	
	// Check length
	if len(alias) < 3 || len(alias) > 50 {
		return fmt.Errorf("alias must be between 3 and 50 characters")
	}
	
	// Check format (alphanumeric, hyphens, underscores only)
	validAlias := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validAlias.MatchString(alias) {
		return fmt.Errorf("alias can only contain letters, numbers, hyphens, and underscores")
	}
	
	// Prevent reserved words
	reserved := []string{"api", "admin", "www", "app", "help", "support", "about"}
	for _, word := range reserved {
		if strings.ToLower(alias) == word {
			return fmt.Errorf("alias '%s' is reserved", alias)
		}
	}
	
	return nil
}

// validateContext checks if the provided context is valid
func validateContext(context string) error {
	if len(context) > 200 {
		return fmt.Errorf("context must be less than 200 characters")
	}
	return nil
}

// ShortenRequest defines the structure for the /api/shorten request body.
type ShortenRequest struct {
	LongURL string `json:"long_url"`
	Alias   string `json:"alias,omitempty"`
	Context string `json:"context,omitempty"`
}

// ShortenResponse defines the structure for the /api/shorten response.
type ShortenResponse struct {
	ShortURL  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// generateShortCode creates a cryptographically secure random string of a fixed length.
func generateShortCode() string {
	b := make([]byte, shortCodeLength)
	for i := range b {
		// Generate a random byte
		randomBytes := make([]byte, 1)
		_, err := rand.Read(randomBytes)
		if err != nil {
			// Fallback to time-based generation if crypto/rand fails
			b[i] = letterBytes[int(time.Now().UnixNano())%len(letterBytes)]
		} else {
			b[i] = letterBytes[int(randomBytes[0])%len(letterBytes)]
		}
	}
	return string(b)
}

func isShortCodeTaken(shortCode string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1)"
	err := db.DB.QueryRow(db.Ctx, query, shortCode).Scan(&exists)
	return exists, err
}

// generateUniqueShortCode creates a random, unique short code.
func generateUniqueShortCode() (string, error) {
	for i := 0; i < 5; i++ { // Retry up to 5 times
		code := generateShortCode()
		taken, err := isShortCodeTaken(code)
		if err != nil {
			return "", err
		}
		if !taken {
			return code, nil
		}
	}
	return "", fmt.Errorf("could not generate a unique short code")
}

// ShortenLink handles the creation of a new shortened link.
func ShortenLink(c *fiber.Ctx) error {
	req := new(ShortenRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate input
	if err := validateURL(req.LongURL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validateAlias(req.Alias); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validateContext(req.Context); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var shortCode string
	var err error

	if req.Alias != "" {
		shortCode = strings.TrimSpace(req.Alias)
		taken, err := isShortCodeTaken(shortCode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
		if taken {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Custom alias is already taken."})
		}
	} else {
		shortCode, err = generateUniqueShortCode()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	expiresAt := time.Now().Add(defaultExpiration)

	// Get user ID from context if available (for authenticated users)
	userID, _ := c.Locals("userID").(string)
	if userID == "" {
		userID = "anonymous" // Default for anonymous users
	}

	// Insert into PostgreSQL
	insertSQL := `INSERT INTO links (short_code, long_url, context, expires_at, user_id) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.DB.Exec(db.Ctx, insertSQL, shortCode, req.LongURL, req.Context, expiresAt, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not save link to database.",
		})
	}

	// Set in Redis cache
	err = db.RDB.Set(db.Ctx, shortCode, req.LongURL, time.Until(expiresAt)).Err()
	if err != nil {
		// Log and ignore cache error
	}

	shortURL := getBaseURL() + "/" + shortCode

	return c.JSON(ShortenResponse{
		ShortURL:  shortURL,
		ExpiresAt: expiresAt,
	})
}

// GenerateQRCode serves a QR code image for a given short link.
func GenerateQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	shortURL := getBaseURL() + "/" + shortCode
	redisKey := "qr:" + shortCode

	// 1. Check if QR code is cached in Redis
	cachedPNG, err := db.RDB.Get(db.Ctx, redisKey).Bytes()
	if err == nil {
		c.Set("Content-Type", "image/png")
		return c.Send(cachedPNG)
	}

	// 2. If not cached, generate a new QR code
	png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate QR code.",
		})
	}

	// 3. Cache the new QR code PNG in Redis
	// We use the original link's expiration duration for the QR code's cache.
	linkExpiresIn, _ := db.RDB.TTL(db.Ctx, shortCode).Result()
	if linkExpiresIn > 0 {
		db.RDB.Set(db.Ctx, redisKey, png, linkExpiresIn).Err()
	}

	c.Set("Content-Type", "image/png")
	return c.Send(png)
}

// RedirectLink handles redirecting a short link to its original URL.
func RedirectLink(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	// 1. Check Redis (cache) first
	longURL, err := db.RDB.Get(db.Ctx, shortCode).Result()
	if err == nil {
		return c.Redirect(longURL, fiber.StatusMovedPermanently)
	}

	// 2. If not in cache, check PostgreSQL
	var expiresAt time.Time
	selectSQL := `SELECT long_url, expires_at FROM links WHERE short_code = $1`
	err = db.DB.QueryRow(db.Ctx, selectSQL, shortCode).Scan(&longURL, &expiresAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Short link not found")
	}

	// 3. Check if the link has expired
	if time.Now().After(expiresAt) {
		return c.Status(fiber.StatusGone).SendString("This link has expired.")
	}

	// 4. Cache the result for future requests
	db.RDB.Set(db.Ctx, shortCode, longURL, time.Until(expiresAt)).Err()

	return c.Redirect(longURL, fiber.StatusMovedPermanently)
} 