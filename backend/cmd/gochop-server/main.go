package main

import (
	"gochop/backend/internal/db"
	"gochop/backend/internal/handlers"
	"gochop/backend/internal/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	// Connect to database and Redis
	if err := db.Connect(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Create the links table if it doesn't exist
	if err := db.CreateLinkTable(); err != nil {
		log.Fatalf("Could not create links table: %v", err)
	}

	// Create the analytics table if it doesn't exist
	if err := db.CreateAnalyticsTable(); err != nil {
		log.Fatalf("Could not create analytics table: %v", err)
	}

	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Add analytics middleware
	app.Use(middleware.AnalyticsMiddleware())

	// Load IP filters from environment
	globalIPFilter := middleware.LoadIPFilterFromEnv()
	adminIPFilter := middleware.CreateAdminIPFilter()

	// Apply global IP filtering if configured
	if globalIPFilter.Mode != "" {
		app.Use(middleware.IPFilterMiddleware(globalIPFilter))
	}

	// Public routes (no authentication required)
	app.Post("/api/shorten", handlers.ShortenLink)
	app.Get("/api/qrcode/:shortCode", handlers.GenerateQRCode)
	app.Get("/:shortCode", handlers.RedirectLink)

	// Authentication routes
	app.Post("/api/auth/login", handlers.Login)
	app.Get("/api/auth/dev-token", handlers.GenerateAdminToken) // Remove in production

	// Protected routes (require authentication)
	protected := app.Group("/api", middleware.JWTMiddleware())
	protected.Get("/profile", handlers.GetProfile)

	// Admin routes (require authentication + admin privileges + optional IP filtering)
	admin := app.Group("/api/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.AdminMiddleware())
	
	// Apply admin-specific IP filtering if configured
	if adminIPFilter.Mode != "" {
		admin.Use(middleware.IPFilterMiddleware(adminIPFilter))
	}
	
	admin.Get("/links", handlers.GetAllLinks)
	admin.Get("/analytics/:shortCode", handlers.GetAnalytics)

	// Start the server
	log.Fatal(app.Listen(":3001")) // Running on port 3001
} 