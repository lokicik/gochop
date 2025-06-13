package main

import (
	"gochop/backend/internal/db"
	"gochop/backend/internal/handlers"
	"gochop/backend/internal/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	// Apply database migrations (runs only pending migrations)
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		// Increase header size limits to prevent "Request Header Fields Too Large" errors
		ReadBufferSize:  32768, // 32KB - increased for NextAuth JWT tokens
		WriteBufferSize: 32768, // 32KB - increased for large responses
		// Allow larger headers for user agents, cookies, etc.
		ServerHeader: "GoChop",
		// Enable better error handling
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Add request logging middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency}) - ${ip} - ${ua}\n",
		TimeFormat: "15:04:05",
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
	app.Get("/api/qrcode/:shortCode", handlers.GenerateQRCode)
	app.Get("/:shortCode", handlers.RedirectLink)

	// Development-only authentication routes have been removed in favor of NextAuth session validation

	// Protected routes (require NextAuth authentication)
	protected := app.Group("/api", middleware.NextAuthMiddleware())
	protected.Get("/profile", handlers.GetProfile)

	// User routes (require authentication)
	user := app.Group("/api/user", middleware.NextAuthMiddleware())
	user.Post("/shorten", handlers.ShortenLink) // Create shortened links (authenticated users only)
	user.Get("/links", handlers.GetAllLinks) // Now returns user's own links or all if admin
	user.Get("/profile", handlers.GetUserProfile) // Full profile with stats
	user.Put("/profile", handlers.UpdateProfile) // Update profile
	user.Get("/stats", handlers.GetUserStats) // User statistics

	// Admin routes (require authentication + admin privileges + optional IP filtering)
	admin := app.Group("/api/admin")
	admin.Use(middleware.NextAuthMiddleware())
	admin.Use(middleware.AdminOnlyMiddleware())
	
	// Apply admin-specific IP filtering if configured
	if adminIPFilter.Mode != "" {
		admin.Use(middleware.IPFilterMiddleware(adminIPFilter))
	}
	
	admin.Get("/links", handlers.GetAllLinks) // Admin can see all links
	admin.Get("/analytics/:shortCode", handlers.GetAnalytics) // Admin can see any link analytics
	admin.Get("/users", handlers.ListUsers) // List all users
	admin.Get("/users/:id", handlers.GetUserByID) // Get specific user details

	// Start the server
	log.Fatal(app.Listen(":3001")) // Running on port 3001
} 