package models

type QuizChoice struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID  uint         `gorm:"not null;index" json:"question_id"`
	Question    QuizQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	ChoiceText  string       `gorm:"type:text;not null" json:"choice_text"`
	ChoiceOrder int          `gorm:"not null" json:"choice_order"`
	IsCorrect   bool         `gorm:"not null;default:false" json:"is_correct"`
}
