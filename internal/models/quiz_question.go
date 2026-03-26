package models

import "time"

type QuizQuestion struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Emoji        string    `gorm:"type:text;not null" json:"emoji"`
	QuestionText string    `gorm:"type:text;not null" json:"question_text"`
	GeneratedAt  time.Time `gorm:"not null;default:now()" json:"generated_at"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`

	Choices []QuizChoice `gorm:"foreignKey:QuestionID" json:"choices,omitempty"`
}
