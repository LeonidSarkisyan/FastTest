package service

import (
	"App/internal/models"
	"errors"
	"github.com/rs/zerolog/log"
)

var (
	TestCreateError = errors.New("ошибка при создании теста")
	TestGetError    = errors.New("ошибка при получении теста")
	TestUpdateError = errors.New("ошибка при обновлении теста")
)

type TestRepository interface {
	Create(test models.Test) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(testID, userID int, title string) error
}

type TestService struct {
	TestRepository
	*QuestionService
}

func NewTestService(r TestRepository, sq *QuestionService) *TestService {
	return &TestService{r, sq}
}

func (s *TestService) Create(title string, userID int) (int, error) {
	newTest := models.Test{
		Title:  title,
		UserID: userID,
	}

	id, err := s.TestRepository.Create(newTest)

	if err != nil {
		log.Err(err).Send()
		return 0, TestCreateError
	}

	_, err = s.QuestionService.Create(id, userID)

	if err != nil {
		log.Err(err).Send()
		return 0, TestCreateError
	}

	return id, nil
}

func (s *TestService) Get(testID, userID int) (models.TestOut, error) {
	test, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return models.TestOut{}, TestGetError
	}

	return test, nil
}

func (s *TestService) GetAll(userID int) ([]models.TestOut, error) {
	tests, err := s.TestRepository.GetAll(userID)

	if err != nil {
		log.Err(err).Send()
		return nil, TestGetError
	}

	return tests, nil
}

func (s *TestService) UpdateTitle(userID, testID int, title string) error {
	err := s.TestRepository.UpdateTitle(testID, userID, title)

	if err != nil {
		log.Err(err).Send()
		return TestUpdateError
	}

	return nil
}
