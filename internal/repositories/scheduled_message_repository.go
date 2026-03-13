package repositories

import (
	"sompong-api/internal/models"

	"gorm.io/gorm"
)

type ScheduledMessageRepository interface {
	Create(msg *models.ScheduledMessage) error
	FindByID(id uint) (*models.ScheduledMessage, error)
	FindByChatID(chatID uint) ([]models.ScheduledMessage, error)
	List() ([]models.ScheduledMessage, error)
	Update(msg *models.ScheduledMessage) error
	Delete(id uint) error
	CreateSubstitutions(subs []models.ScheduledMessageSubstitution) error
	DeleteSubstitutionsByMessageID(id uint) error
}

type scheduledMessageRepository struct {
	db *gorm.DB
}

func NewScheduledMessageRepository(db *gorm.DB) ScheduledMessageRepository {
	return &scheduledMessageRepository{db: db}
}

func (r *scheduledMessageRepository) Create(msg *models.ScheduledMessage) error {
	return r.db.Create(msg).Error
}

func (r *scheduledMessageRepository) FindByID(id uint) (*models.ScheduledMessage, error) {
	var msg models.ScheduledMessage
	err := r.db.
		Preload("Substitutions.User").
		Preload("Chat").
		Preload("CreatedByUser").
		First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *scheduledMessageRepository) FindByChatID(chatID uint) ([]models.ScheduledMessage, error) {
	var msgs []models.ScheduledMessage
	err := r.db.
		Preload("Substitutions.User").
		Where("chat_id = ?", chatID).
		Find(&msgs).Error
	return msgs, err
}

func (r *scheduledMessageRepository) List() ([]models.ScheduledMessage, error) {
	var msgs []models.ScheduledMessage
	err := r.db.
		Preload("Substitutions.User").
		Find(&msgs).Error
	return msgs, err
}

func (r *scheduledMessageRepository) Update(msg *models.ScheduledMessage) error {
	return r.db.Save(msg).Error
}

func (r *scheduledMessageRepository) Delete(id uint) error {
	return r.db.Delete(&models.ScheduledMessage{}, id).Error
}

func (r *scheduledMessageRepository) CreateSubstitutions(subs []models.ScheduledMessageSubstitution) error {
	if len(subs) == 0 {
		return nil
	}
	return r.db.Create(&subs).Error
}

func (r *scheduledMessageRepository) DeleteSubstitutionsByMessageID(id uint) error {
	return r.db.Where("scheduled_message_id = ?", id).Delete(&models.ScheduledMessageSubstitution{}).Error
}
