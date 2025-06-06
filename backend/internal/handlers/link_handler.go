package handlers

import (
	"fmt"
	"gochop/backend/internal/db"
	"math/rand"
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

// ShortenRequest defines the structure for the /api/shorten request body.
type ShortenRequest struct {
	LongURL string `json:"long_url"`
	Alias   string `json:"alias,omitempty"`
}

// ShortenResponse defines the structure for the /api/shorten response.
type ShortenResponse struct {
	ShortURL  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// init seeds the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateShortCode creates a random string of a fixed length.
func generateShortCode() string {
	b := make([]byte, shortCodeLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
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

	// Insert into PostgreSQL
	insertSQL := `INSERT INTO links (short_code, long_url, expires_at) VALUES ($1, $2, $3)`
	_, err = db.DB.Exec(db.Ctx, insertSQL, shortCode, req.LongURL, expiresAt)
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

	shortURL := "http://localhost:3001/" + shortCode

	return c.JSON(ShortenResponse{
		ShortURL:  shortURL,
		ExpiresAt: expiresAt,
	})
}

// GenerateQRCode serves a QR code image for a given short link.
func GenerateQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	shortURL := "http://localhost:3001/" + shortCode
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