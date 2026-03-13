package repositories

import (
	"sompong-api/internal/models"

	"gorm.io/gorm"
)

type ChatRepository interface {
	List() ([]models.Chat, error)
	FindUsersByChatID(chatID uint) ([]models.User, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) List() ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Preload("Group").Find(&chats).Error
	return chats, err
}

func (r *chatRepository) FindUsersByChatID(chatID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.
		Distinct("users.*").
		Joins("JOIN chat_messages ON chat_messages.user_id = users.id").
		Where("chat_messages.chat_id = ? AND chat_messages.user_id IS NOT NULL", chatID).
		Find(&users).Error
	return users, err
}
