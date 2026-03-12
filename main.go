package main

import (
	"log"
	"sompong-api/internal/config"
	"sompong-api/internal/database"
	"sompong-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Load()
	database.Connect()

	app := fiber.New(fiber.Config{
		AppName: "sompong-api",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	routes.Setup(app)

	log.Fatal(app.Listen(":" + config.App.AppPort))
}
