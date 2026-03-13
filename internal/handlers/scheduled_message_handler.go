package handlers

import (
	"strconv"

	"sompong-api/internal/dto"
	"sompong-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ScheduledMessageHandler struct {
	service services.ScheduledMessageService
}

func NewScheduledMessageHandler(service services.ScheduledMessageService) *ScheduledMessageHandler {
	return &ScheduledMessageHandler{service: service}
}

func (h *ScheduledMessageHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateScheduledMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	msg, err := h.service.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(msg)
}

func (h *ScheduledMessageHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	msg, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(msg)
}

func (h *ScheduledMessageHandler) List(c *fiber.Ctx) error {
	chatIDStr := c.Query("chat_id")
	if chatIDStr == "" {
		msgs, err := h.service.List()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(msgs)
	}

	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat_id"})
	}

	msgs, err := h.service.ListByChatID(uint(chatID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(msgs)
}

func (h *ScheduledMessageHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req dto.UpdateScheduledMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	msg, err := h.service.Update(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(msg)
}

func (h *ScheduledMessageHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
