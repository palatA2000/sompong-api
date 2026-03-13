package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DisplayName *string   `gorm:"type:varchar" json:"display_name"`
	LineUserID  string    `gorm:"type:varchar;not null;unique" json:"line_user_id"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}
