package models

import "time"

type ScheduledMessageSubstitution struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ScheduledMessageID uint      `gorm:"not null;index" json:"scheduled_message_id"`
	ScheduledMessage   ScheduledMessage `gorm:"foreignKey:ScheduledMessageID" json:"scheduled_message,omitempty"`
	Key                string    `gorm:"type:varchar;not null" json:"key"`
	Type               string    `gorm:"type:varchar;not null" json:"type"`
	MentioneeType      *string   `gorm:"type:varchar" json:"mentionee_type"`
	UserID             *uint     `gorm:"index" json:"user_id"`
	User               *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt          time.Time `gorm:"not null;default:now()" json:"created_at"`
}
