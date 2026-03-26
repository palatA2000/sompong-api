package repositories

import (
	"sompong-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuizRepository interface {
	GetQuestionByID(id uint) (*models.QuizQuestion, error)
	CreateQuestion(question *models.QuizQuestion) error
	CreateQuestionWithChoices(question *models.QuizQuestion, choices []models.QuizChoice) error
	CreateChoices(choices []models.QuizChoice) error
	GetChoicesByQuestionID(questionID uint) ([]models.QuizChoice, error)
	GetChoiceByID(id uint) (*models.QuizChoice, error)
	UpsertUserByLineID(lineUserID string) (*models.User, error)
	UpsertGroupByLineID(lineGroupID string) (*models.Group, error)
	CreateAttemptIfNotExists(attempt *models.QuizAttempt) (bool, error)
	IncrementUserScore(userID uint, delta int) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) GetQuestionByID(id uint) (*models.QuizQuestion, error) {
	var q models.QuizQuestion
	err := r.db.Preload("Choices").First(&q, id).Error
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *quizRepository) CreateQuestion(question *models.QuizQuestion) error {
	return r.db.Create(question).Error
}

func (r *quizRepository) CreateQuestionWithChoices(question *models.QuizQuestion, choices []models.QuizChoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(question).Error; err != nil {
			return err
		}
		if len(choices) == 0 {
			return nil
		}
		for i := range choices {
			choices[i].QuestionID = question.ID
		}
		return tx.Create(&choices).Error
	})
}

func (r *quizRepository) CreateChoices(choices []models.QuizChoice) error {
	if len(choices) == 0 {
		return nil
	}
	return r.db.Create(&choices).Error
}

func (r *quizRepository) GetChoicesByQuestionID(questionID uint) ([]models.QuizChoice, error) {
	var choices []models.QuizChoice
	err := r.db.Where("question_id = ?", questionID).Order("choice_order asc").Find(&choices).Error
	return choices, err
}

func (r *quizRepository) GetChoiceByID(id uint) (*models.QuizChoice, error) {
	var choice models.QuizChoice
	err := r.db.First(&choice, id).Error
	if err != nil {
		return nil, err
	}
	return &choice, nil
}

func (r *quizRepository) UpsertUserByLineID(lineUserID string) (*models.User, error) {
	user := &models.User{LineUserID: lineUserID}
	if err := r.db.FirstOrCreate(user, models.User{LineUserID: lineUserID}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *quizRepository) UpsertGroupByLineID(lineGroupID string) (*models.Group, error) {
	group := &models.Group{LineGroupID: lineGroupID}
	if err := r.db.FirstOrCreate(group, models.Group{LineGroupID: lineGroupID}).Error; err != nil {
		return nil, err
	}
	return group, nil
}

func (r *quizRepository) CreateAttemptIfNotExists(attempt *models.QuizAttempt) (bool, error) {
	res := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(attempt)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *quizRepository) IncrementUserScore(userID uint, delta int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("score", gorm.Expr("score + ?", delta)).Error
}
