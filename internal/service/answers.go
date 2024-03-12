package service

import (
	"App/internal/models"
	"errors"
)

var (
	AnswerCreateError = errors.New("ошибка при создании ответа")
	AnswerUpdateError = errors.New("ошибка при обновлении ответа")
	AnswerDeleteError = errors.New("ошибка при удалении ответа")
)

type AnswerRepository interface {
	Create(questionID int) (int, error)
	CreateThree(questionID int) error
	GetAll(questionID int) ([]models.Answer, error)
	Update(answerUpdate models.AnswerUpdate, answerID, questionID int) error
	Delete(answerID, questionID int) error
}

type AnswerService struct {
	AnswerRepository
	TestRepository
	QuestionRepository
}

func NewAnswerService(a AnswerRepository, t TestRepository, q QuestionRepository) *AnswerService {
	return &AnswerService{a, t, q}
}

func (s *AnswerService) Create(userID, testID, questionID int) (int, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		return 0, TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		return 0, QuestionNotFound
	}

	id, err := s.AnswerRepository.Create(questionID)

	if err != nil {
		return 0, AnswerCreateError
	}

	return id, nil
}

func (s *AnswerService) CreateThree(userID, testID, questionID int) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		return TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		return QuestionNotFound
	}

	err = s.AnswerRepository.CreateThree(questionID)

	if err != nil {
		return AnswerCreateError
	}

	return nil
}

func (s *AnswerService) GetAllByQuestionID(userID, testID, questionID int) ([]models.Answer, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		return nil, TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		return nil, QuestionNotFound
	}

	answers, err := s.AnswerRepository.GetAll(questionID)

	if err != nil {
		return nil, AnswerCreateError
	}

	return answers, nil
}

func (s *AnswerService) Update(userID, testID, questionID, answerID int, answerUpdate models.AnswerUpdate) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		return TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		return QuestionNotFound
	}

	err = s.AnswerRepository.Update(answerUpdate, answerID, questionID)

	if err != nil {
		return AnswerUpdateError
	}

	return nil
}

func (s *AnswerService) Delete(userID, testID, questionID, answerID int) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		return TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		return QuestionNotFound
	}

	err = s.AnswerRepository.Delete(answerID, questionID)

	if err != nil {
		return AnswerDeleteError
	}

	return nil
}
