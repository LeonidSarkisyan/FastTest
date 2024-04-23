package service

import (
	"App/internal/models"
	"App/internal/questions"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	QuestionErrorCreate = errors.New("ошибка при создании вопроса")
	QuestionGetCreate   = errors.New("ошибка при получении вопросов")
	QuestionUpdateError = errors.New("ошибка при обновлении вопроса")
	QuestionDeleteError = errors.New("ошибка при удалении вопроса")
	QuestionNotFound    = errors.New("вопрос не найден")

	TestForbidden = errors.New("тест не найден")

	NotFountQuestionType = errors.New("неизвестный тип вопроса")

	InCorrectStructureError = errors.New("некорректная структура")
	ImageErrorCreate        = errors.New("ошибка при сохранении файла")
	ImageDeleteCreate       = errors.New("ошибка при удалении файла")
)

const (
	MediaImagePath = "/static/media/questions/"

	Group = "group"
)

type QuestionRepository interface {
	Create(testID int) (int, error)
	GetAll(testID int) ([]models.Question, error)
	Update(questionID, testID int, question models.QuestionUpdate) error
	Get(questionID, testID int) (models.Question, error)
	Delete(questionID, testID int) error
	CreateWithType(testID int, type_ string, data []byte) (int, error)
	Save(testID, questionID int, type_ string, data []byte) error

	UploadImage(userID, testID, questionID int, filename string) error
	DeleteImage(questionID int) ([]string, error)

	GetAllWithAnswers(testID int) ([]models.QuestionWithAnswers, error)
	CreateManyQuestions(
		testID int, questions []models.QuestionWithAnswersWithOutIsCorrect,
	) ([]models.QuestionWithAnswersWithOutIsCorrect, error)
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

func (s *QuestionService) CreateWithType(testID, userID int, type_ string) (int, any, error) {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return 0, make([]byte, 0), TestForbidden
	}

	data := questions.CreateDataForType(type_)

	id, err := s.QuestionRepository.CreateWithType(testID, type_, data)

	if err != nil {
		log.Err(err).Send()
		return 0, make([]byte, 0), TestForbidden
	}

	var d any

	err = json.Unmarshal(data, &d)

	if err != nil {
		log.Err(err).Send()
		return 0, make([]byte, 0), TestForbidden
	}

	return id, d, nil
}

func (s *QuestionService) Save(testID, userID, questionID int, type_ string, data []byte) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return TestForbidden
	}

	log.Info().Any("DATA", data).Send()

	err = questions.CheckJSONStructure(type_, data)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	err = s.QuestionRepository.Save(testID, questionID, type_, data)

	if err != nil {
		log.Err(err).Send()
		return InCorrectStructureError
	}

	return nil
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

func (s *QuestionService) UploadImage(userID, testID, questionID int, filename string) (string, error) {
	err := s.DeleteImage(userID, testID, questionID)

	if err != nil {
		log.Err(err).Send()
		return "", ImageDeleteCreate
	}

	err = s.QuestionRepository.UploadImage(userID, testID, questionID, MediaImagePath+filename)

	if err != nil {
		return "", ImageErrorCreate
	}

	return MediaImagePath + filename, nil
}

func (s *QuestionService) DeleteImage(userID, testID, questionID int) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return TestForbidden
	}

	_, err = s.QuestionRepository.Get(questionID, testID)

	if err != nil {
		log.Err(err).Send()
		return QuestionGetCreate
	}

	urls, err := s.QuestionRepository.DeleteImage(questionID)

	if err != nil {
		log.Err(err).Send()
		return ImageDeleteCreate
	}

	log.Info().Strs("urls", urls).Send()

	for _, url := range urls {
		err = os.Remove(url)
		if err != nil {
			log.Err(err).Send()
		}
	}

	return nil
}
