package routes

import (
	"sompong-api/internal/database"
	"sompong-api/internal/handlers"
	"sompong-api/internal/middleware"
	"sompong-api/internal/repositories"
	"sompong-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Quiz routes (API key auth)
	quizRepo := repositories.NewQuizRepository(database.DB)
	quizService := services.NewQuizService(quizRepo)
	quizHandler := handlers.NewQuizHandler(quizService)

	quiz := api.Group("/quiz", middleware.APIKeyAuth())
	quiz.Post("/", quizHandler.GenerateQuiz)
	quiz.Post("/answers", quizHandler.Answer)

	// Protected routes (require LINE access token)
	// protected := api.Group("/", middleware.LIFFAuth())

	// // Scheduled Messages
	// smRepo := repositories.NewScheduledMessageRepository(database.DB)
	// smService := services.NewScheduledMessageService(smRepo)
	// smHandler := handlers.NewScheduledMessageHandler(smService)

	// sm := protected.Group("/scheduled-messages")
	// sm.Post("/", smHandler.Create)
	// sm.Get("/", smHandler.List)
	// sm.Get("/:id", smHandler.GetByID)
	// sm.Put("/:id", smHandler.Update)
	// sm.Delete("/:id", smHandler.Delete)

	// // Chats
	// chatRepo := repositories.NewChatRepository(database.DB)
	// chatHandler := handlers.NewChatHandler(chatRepo)

	// chats := protected.Group("/chats")
	// chats.Get("/", chatHandler.List)
	// chats.Get("/:id/users", chatHandler.ListUsers)
}
