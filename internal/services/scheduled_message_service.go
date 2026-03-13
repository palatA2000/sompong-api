package services

import (
	"fmt"
	"sompong-api/internal/dto"
	"sompong-api/internal/models"
	"sompong-api/internal/repositories"
)

type ScheduledMessageService interface {
	Create(req dto.CreateScheduledMessageRequest) (*models.ScheduledMessage, error)
	GetByID(id uint) (*models.ScheduledMessage, error)
	ListByChatID(chatID uint) ([]models.ScheduledMessage, error)
	List() ([]models.ScheduledMessage, error)
	Update(id uint, req dto.UpdateScheduledMessageRequest) (*models.ScheduledMessage, error)
	Delete(id uint) error
}

type scheduledMessageService struct {
	repo repositories.ScheduledMessageRepository
}

func NewScheduledMessageService(repo repositories.ScheduledMessageRepository) ScheduledMessageService {
	return &scheduledMessageService{repo: repo}
}

func (s *scheduledMessageService) Create(req dto.CreateScheduledMessageRequest) (*models.ScheduledMessage, error) {
	msg := &models.ScheduledMessage{
		ChatID:            req.ChatID,
		CreatedByUserID:   req.CreatedByUserID,
		MessageText:       req.MessageText,
		MessageType:       req.MessageType,
		Timezone:          req.Timezone,
		StartAt:           req.StartAt,
		EndAt:             req.EndAt,
		FrequencyType:     req.FrequencyType,
		FrequencyInterval: req.FrequencyInterval,
		DaysOfWeek:        req.DaysOfWeek,
		DayOfMonth:        req.DayOfMonth,
		TimeOfDay:         req.TimeOfDay,
		CronExpression:    req.CronExpression,
		IsActive:          req.IsActive,
	}

	if err := s.repo.Create(msg); err != nil {
		return nil, err
	}

	if err := s.repo.CreateSubstitutions(buildSubstitutions(msg.ID, req.MentionUserIDs)); err != nil {
		return nil, err
	}

	return s.repo.FindByID(msg.ID)
}

func (s *scheduledMessageService) GetByID(id uint) (*models.ScheduledMessage, error) {
	return s.repo.FindByID(id)
}

func (s *scheduledMessageService) ListByChatID(chatID uint) ([]models.ScheduledMessage, error) {
	return s.repo.FindByChatID(chatID)
}

func (s *scheduledMessageService) List() ([]models.ScheduledMessage, error) {
	return s.repo.List()
}

func (s *scheduledMessageService) Update(id uint, req dto.UpdateScheduledMessageRequest) (*models.ScheduledMessage, error) {
	msg, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.MessageText != nil {
		msg.MessageText = *req.MessageText
	}
	if req.MessageType != nil {
		msg.MessageType = *req.MessageType
	}
	if req.Timezone != nil {
		msg.Timezone = *req.Timezone
	}
	if req.StartAt != nil {
		msg.StartAt = *req.StartAt
	}
	if req.EndAt != nil {
		msg.EndAt = req.EndAt
	}
	if req.FrequencyType != nil {
		msg.FrequencyType = *req.FrequencyType
	}
	if req.FrequencyInterval != nil {
		msg.FrequencyInterval = req.FrequencyInterval
	}
	if req.DaysOfWeek != nil {
		msg.DaysOfWeek = req.DaysOfWeek
	}
	if req.DayOfMonth != nil {
		msg.DayOfMonth = req.DayOfMonth
	}
	if req.TimeOfDay != nil {
		msg.TimeOfDay = req.TimeOfDay
	}
	if req.CronExpression != nil {
		msg.CronExpression = req.CronExpression
	}
	if req.IsActive != nil {
		msg.IsActive = *req.IsActive
	}

	if err := s.repo.Update(msg); err != nil {
		return nil, err
	}

	if req.MentionUserIDs != nil {
		if err := s.repo.DeleteSubstitutionsByMessageID(id); err != nil {
			return nil, err
		}
		if err := s.repo.CreateSubstitutions(buildSubstitutions(id, req.MentionUserIDs)); err != nil {
			return nil, err
		}
	}

	return s.repo.FindByID(id)
}

func (s *scheduledMessageService) Delete(id uint) error {
	if err := s.repo.DeleteSubstitutionsByMessageID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func buildSubstitutions(msgID uint, userIDs []uint) []models.ScheduledMessageSubstitution {
	subs := make([]models.ScheduledMessageSubstitution, 0, len(userIDs))
	mentioneeType := "user"
	for i, uid := range userIDs {
		subs = append(subs, models.ScheduledMessageSubstitution{
			ScheduledMessageID: msgID,
			Key:                fmt.Sprintf("USER_%d", i+1),
			Type:               "mention",
			MentioneeType:      &mentioneeType,
			UserID:             &uid,
		})
	}
	return subs
}
