package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/nowel-xyz/quiz/database"
	"github.com/nowel-xyz/quiz/routers"
)

func main() {
	app := fiber.New()

	err := database.Init("mongodb://localhost:27017", "kahoot")
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitRedis()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}


	// Register all routers
	app.Static("/static", "./public")
	routers.SetupQuizRoutes(app)
	routers.SetupAuthRoutes(app)
	routers.SetupWebSocket(app)
	routers.SetupHomeRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("SECRET_KEY from env:", os.Getenv("SECRET_KEY"))
}


