package models

import "time"

type Chat struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID string    `gorm:"type:varchar;not null;unique" json:"conversation_id"`
	SourceType     string    `gorm:"type:varchar;not null" json:"source_type"`
	SourceID       *string   `gorm:"type:varchar" json:"source_id"`
	GroupID        *uint     `gorm:"index" json:"group_id"`
	Group          *Group    `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	LastMessageAt  time.Time `gorm:"not null;default:now()" json:"last_message_at"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:now()" json:"updated_at"`
}
