package handlers

import (
	"strconv"

	"sompong-api/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	repo repositories.ChatRepository
}

func NewChatHandler(repo repositories.ChatRepository) *ChatHandler {
	return &ChatHandler{repo: repo}
}

func (h *ChatHandler) List(c *fiber.Ctx) error {
	chats, err := h.repo.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chats)
}

func (h *ChatHandler) ListUsers(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	users, err := h.repo.FindUsersByChatID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}
