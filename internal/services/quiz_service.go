package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sompong-api/internal/config"
	"sompong-api/internal/dto"
	"sompong-api/internal/models"
	"sompong-api/internal/repositories"
)

type QuizService interface {
	GenerateQuiz(ctx context.Context) (*models.QuizQuestion, error)
	SubmitAnswer(ctx context.Context, req dto.QuizAnswerRequest) (*dto.QuizAnswerResponse, error)
}

type quizService struct {
	repo   repositories.QuizRepository
	client *http.Client
}

func NewQuizService(repo repositories.QuizRepository) QuizService {
	return &quizService{repo: repo, client: &http.Client{Timeout: 15 * time.Second}}
}

func (s *quizService) GenerateQuiz(ctx context.Context) (*models.QuizQuestion, error) {
	question, err := s.generateAndStoreQuestion(ctx)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (s *quizService) SubmitAnswer(ctx context.Context, req dto.QuizAnswerRequest) (*dto.QuizAnswerResponse, error) {
	_ = ctx
	if strings.TrimSpace(req.LineUserID) == "" {
		return nil, fmt.Errorf("line_user_id is required")
	}
	if req.QuestionID == 0 || req.ChoiceID == 0 {
		return nil, fmt.Errorf("question_id and choice_id are required")
	}

	user, err := s.repo.UpsertUserByLineID(req.LineUserID)
	if err != nil {
		return nil, err
	}

	choice, err := s.repo.GetChoiceByID(req.ChoiceID)
	if err != nil {
		return nil, err
	}
	if choice.QuestionID != req.QuestionID {
		return nil, fmt.Errorf("choice does not belong to question")
	}

	attempt := &models.QuizAttempt{
		QuestionID: req.QuestionID,
		UserID:     user.ID,
		ChoiceID:   req.ChoiceID,
		IsCorrect:  choice.IsCorrect,
		AnsweredAt: time.Now().UTC(),
	}

	created, err := s.repo.CreateAttemptIfNotExists(attempt)
	if err != nil {
		return nil, err
	}
	if !created {
		return &dto.QuizAnswerResponse{
			Correct:         false,
			AlreadyAnswered: true,
			Score:           user.Score,
		}, nil
	}

	if choice.IsCorrect {
		if err := s.repo.IncrementUserScore(user.ID, 1); err != nil {
			return nil, err
		}
		user.Score += 1
	}

	return &dto.QuizAnswerResponse{
		Correct:         choice.IsCorrect,
		AlreadyAnswered: false,
		Score:           user.Score,
	}, nil
}

type geminiQuizPayload struct {
	Emoji        string   `json:"emoji"`
	Question     string   `json:"question"`
	Choices      []string `json:"choices"`
	CorrectIndex int      `json:"correct_index"`
}

func (s *quizService) generateAndStoreQuestion(ctx context.Context) (*models.QuizQuestion, error) {
	payload, err := s.generateQuestionWithGemini(ctx, config.App.QuizChoiceCount)
	if err != nil {
		return nil, err
	}
	question := &models.QuizQuestion{
		Emoji:        payload.Emoji,
		QuestionText: payload.Question,
		GeneratedAt:  time.Now().UTC(),
	}

	choices := make([]models.QuizChoice, 0, len(payload.Choices))
	for i, c := range payload.Choices {
		choices = append(choices, models.QuizChoice{
			ChoiceText:  c,
			ChoiceOrder: i + 1,
			IsCorrect:   i == payload.CorrectIndex,
		})
	}

	if err := s.repo.CreateQuestionWithChoices(question, choices); err != nil {
		return nil, err
	}
	return s.repo.GetQuestionByID(question.ID)
}

func (s *quizService) generateQuestionWithGemini(ctx context.Context, choiceCount int) (*geminiQuizPayload, error) {
	if config.App.GeminiAPIKey == "" {
		return nil, fmt.Errorf("missing GEMINI_API_KEY")
	}
	if choiceCount <= 1 {
		choiceCount = 4
	}
	prompt := buildGeminiPrompt(choiceCount)

	body := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": prompt},
				},
			},
		},
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", strings.TrimRight(config.App.GeminiBaseURL, "/"), config.App.Gemini3Model, config.App.GeminiAPIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var gr geminiResponse
		json.NewDecoder(resp.Body).Decode(&gr)
		fmt.Println(gr, "gr")
		return nil, fmt.Errorf("gemini response status: %d", resp.StatusCode)
	}

	var gr geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}

	text := gr.FirstText()
	if text == "" {
		return nil, fmt.Errorf("empty gemini response")
	}

	jsonText := extractJSON(text)
	if jsonText == "" {
		return nil, fmt.Errorf("invalid gemini response format")
	}

	var payload geminiQuizPayload
	if err := json.Unmarshal([]byte(jsonText), &payload); err != nil {
		return nil, err
	}
	if err := validateGeminiPayload(&payload, choiceCount); err != nil {
		return nil, err
	}
	return &payload, nil
}

func buildGeminiPrompt(choiceCount int) string {
	return fmt.Sprintf(
		"You are a professional Thai music trivia engine. Create one challenge for a Thai song emoji guessing game."+
			"Select a highly recognizable Thai song (Pop, Rock, or Indie)."+
			"Convert the song title into 2-5 emojis that logically represent the words or the literal meaning of the title."+
			"Generate 4 choices in Thai. The 3 distractors must be real Thai song titles that are within the same genre or era as the correct one to ensure difficulty."+
			"Ensure 'correct_index' matches the position of the correct answer in the 'choices' array (0-based)."+
			"Output Format: Return ONLY a strictly valid JSON object. No markdown, no prose."+
			"Schema: {\"emoji\": string, \"question\": string, \"choices\": string[%d], \"correct_index\": number}. "+
			"Constraints: The 'question' field must always be 'เพลงนี้คือเพลงอะไร?'. 'choices' must contain exactly 4 Thai song titles. Emojis must be clear and not overly obscure.",
		choiceCount,
	)
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r geminiResponse) FirstText() string {
	if len(r.Candidates) == 0 {
		return ""
	}
	if len(r.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	return r.Candidates[0].Content.Parts[0].Text
}

func extractJSON(text string) string {
	trimmed := strings.TrimSpace(text)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return trimmed[start : end+1]
}

func validateGeminiPayload(p *geminiQuizPayload, choiceCount int) error {
	if strings.TrimSpace(p.Emoji) == "" {
		return fmt.Errorf("emoji is required")
	}
	if strings.TrimSpace(p.Question) == "" {
		return fmt.Errorf("question is required")
	}
	if len(p.Choices) != choiceCount {
		return fmt.Errorf("choices must have %d items", choiceCount)
	}
	for i := range p.Choices {
		if strings.TrimSpace(p.Choices[i]) == "" {
			return fmt.Errorf("choice %d is empty", i)
		}
	}
	if p.CorrectIndex < 0 || p.CorrectIndex >= len(p.Choices) {
		return fmt.Errorf("correct_index out of range")
	}
	return nil
}
