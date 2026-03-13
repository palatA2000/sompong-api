package models

import "time"

type ChatMessage struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID uint   `gorm:"not null;index" json:"chat_id"`
	Chat   Chat   `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	UserID *uint  `gorm:"index" json:"user_id"`
	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Text   string `gorm:"type:text;not null" json:"text"`
	SentAt time.Time `gorm:"not null" json:"sent_at"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
}
