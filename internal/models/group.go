package models

import "time"

type Group struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LineGroupID string    `gorm:"type:varchar;not null;unique" json:"line_group_id"`
	GroupName   *string   `gorm:"type:varchar" json:"group_name"`
	PictureURL  *string   `gorm:"type:text" json:"picture_url"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}
