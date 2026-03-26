package models

import "time"

type QuizAttempt struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID uint      `gorm:"not null;index:idx_quiz_attempt_unique,unique" json:"question_id"`
	UserID     uint      `gorm:"not null;index:idx_quiz_attempt_unique,unique" json:"user_id"`
	ChoiceID   uint      `gorm:"not null" json:"choice_id"`
	IsCorrect  bool      `gorm:"not null" json:"is_correct"`
	AnsweredAt time.Time `gorm:"not null" json:"answered_at"`
}
