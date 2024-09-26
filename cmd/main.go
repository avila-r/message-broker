package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Printf("failed while loading .env file - %v", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	url := os.Getenv("SERVER_URL")

	app.Listen(url)
}
