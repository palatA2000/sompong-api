package routes

import (
	"sompong-api/internal/database"
	"sompong-api/internal/handlers"
	"sompong-api/internal/repositories"
	"sompong-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Scheduled Messages
	smRepo := repositories.NewScheduledMessageRepository(database.DB)
	smService := services.NewScheduledMessageService(smRepo)
	smHandler := handlers.NewScheduledMessageHandler(smService)

	sm := api.Group("/scheduled-messages")
	sm.Post("/", smHandler.Create)
	sm.Get("/", smHandler.List)
	sm.Get("/:id", smHandler.GetByID)
	sm.Put("/:id", smHandler.Update)
	sm.Delete("/:id", smHandler.Delete)

	// Chats
	chatRepo := repositories.NewChatRepository(database.DB)
	chatHandler := handlers.NewChatHandler(chatRepo)

	chats := api.Group("/chats")
	chats.Get("/", chatHandler.List)
	chats.Get("/:id/users", chatHandler.ListUsers)
}
