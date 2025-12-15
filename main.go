package main

import (
	"log"
	"os"

	"inventory-api/config"
	"inventory-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: File .env tidak ditemukan")
	}

	// 2. Connect Database
	config.ConnectDB()

	// 3. Init Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// 4. Setup Routes
	routes.SetupRoutes(app)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}