package main

import (
	"gochop/backend/internal/db"
	"gochop/backend/internal/handlers"
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

	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Define routes
	app.Post("/api/shorten", handlers.ShortenLink)
	app.Get("/api/qrcode/:shortCode", handlers.GenerateQRCode)
	app.Get("/:shortCode", handlers.RedirectLink)

	// Start the server
	log.Fatal(app.Listen(":3001")) // Running on port 3001
} 