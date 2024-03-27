package service

import (
	"App/internal/models"
	"errors"
	"github.com/rs/zerolog/log"
)

var (
	QuestionErrorCreate = errors.New("ошибка при создании вопроса")
	QuestionGetCreate   = errors.New("ошибка при получении вопросов")
	QuestionUpdateError = errors.New("ошибка при обновлении вопроса")
	QuestionDeleteError = errors.New("ошибка при удалении вопроса")
	QuestionNotFound    = errors.New("вопрос не найден")

	TestForbidden = errors.New("тест не найден")
)

type QuestionRepository interface {
	Create(testID int) (int, error)
	GetAll(testID int) ([]models.Question, error)
	Update(questionID, testID int, question models.QuestionUpdate) error
	Get(questionID, testID int) (models.Question, error)
	Delete(questionID, testID int) error

	GetAllWithAnswers(testID int) ([]models.QuestionWithAnswers, error)
}

type QuestionService struct {
	QuestionRepository
	TestRepository
	*AnswerService
}

func NewQuestionService(r QuestionRepository, rt TestRepository, sa *AnswerService) *QuestionService {
	return &QuestionService{
		QuestionRepository: r,
		TestRepository:     rt,
		AnswerService:      sa,
	}
}

func (s *QuestionService) Create(testID, userID int) (int, []int, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return 0, nil, TestForbidden
	}

	id, err := s.QuestionRepository.Create(testID)

	if err != nil {
		log.Err(err).Send()
		return 0, nil, QuestionErrorCreate
	}

	ids, err := s.AnswerService.CreateThree(userID, testID, id)

	if err != nil {
		return 0, nil, AnswerCreateError
	}

	return id, ids, nil
}

func (s *QuestionService) GetAll(testID, userID int) ([]models.Question, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return nil, TestForbidden
	}

	questions, err := s.QuestionRepository.GetAll(testID)

	if err != nil {
		log.Err(err).Send()
		return nil, QuestionGetCreate
	}

	return questions, nil
}

func (s *QuestionService) Update(userID, testID, questionID int, question models.QuestionUpdate) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return TestForbidden
	}

	err = s.QuestionRepository.Update(questionID, testID, question)

	if err != nil {
		log.Err(err).Send()
		return QuestionUpdateError
	}

	return nil
}

func (s *QuestionService) Delete(userID, testID, questionID int) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return TestForbidden
	}

	err = s.QuestionRepository.Delete(questionID, testID)

	if err != nil {
		log.Err(err).Send()
		return QuestionDeleteError
	}

	return nil
}

func (s *QuestionService) GetAllQuestionsWithAnswers(testID int) ([]models.QuestionWithAnswers, error) {
	questions, err := s.QuestionRepository.GetAllWithAnswers(testID)

	if err != nil {
		return nil, err
	}

	return questions, nil
}
