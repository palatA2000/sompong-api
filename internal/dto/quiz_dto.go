package dto

type QuizAnswerRequest struct {
	LineUserID string `json:"line_user_id"`
	QuestionID uint   `json:"question_id"`
	ChoiceID   uint   `json:"choice_id"`
}

type QuizAnswerResponse struct {
	Correct         bool `json:"correct"`
	AlreadyAnswered bool `json:"already_answered"`
	Score           int  `json:"score"`
}
