package handlers

import (
	"sompong-api/internal/dto"
	"sompong-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type QuizHandler struct {
	service services.QuizService
}

func NewQuizHandler(service services.QuizService) *QuizHandler {
	return &QuizHandler{service: service}
}

func (h *QuizHandler) GenerateQuiz(c *fiber.Ctx) error {
	resp, err := h.service.GenerateQuiz(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

func (h *QuizHandler) Answer(c *fiber.Ctx) error {
	var req dto.QuizAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	resp, err := h.service.SubmitAnswer(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
