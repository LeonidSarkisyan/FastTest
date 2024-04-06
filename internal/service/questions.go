package service

import (
	"App/internal/models"
	"encoding/json"
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

	NotFountQuestionType = errors.New("неизвестный тип вопроса")

	InCorrectStructureError = errors.New("некорректная структура")
)

const (
	Group = "group"
	Range = "range"

	CountGroups        = 2
	CountAnswerInGroup = 1
	CountRanges        = 3
)

type QuestionRepository interface {
	Create(testID int) (int, error)
	GetAll(testID int) ([]models.Question, error)
	Update(questionID, testID int, question models.QuestionUpdate) error
	Get(questionID, testID int) (models.Question, error)
	Delete(questionID, testID int) error
	CreateWithType(testID int, type_ string, data []byte) (int, error)
	Save(testID, questionID int, type_ string, data []byte) error

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

	switch type_ {
	case Group:
		data := models.QuestionGroupData{Groups: make([]models.Group, CountGroups)}

		for i := range CountGroups {
			data.Groups[i].Name = ""
			data.Groups[i].Answers = make([]string, CountAnswerInGroup)
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Err(err).Send()
			return 0, make([]byte, 0), InCorrectStructureError
		}

		id, err := s.QuestionRepository.CreateWithType(testID, type_, jsonData)

		if err != nil {
			log.Err(err).Send()
			return 0, make([]byte, 0), TestForbidden
		}

		return id, data.Groups, nil
	case Range:
		data := models.QuestionRangeData{Ranges: make([]models.Range, CountRanges)}

		for i := range CountRanges {
			data.Ranges[i].Text = ""
			data.Ranges[i].Index = i
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Err(err).Send()
			return 0, make([]byte, 0), InCorrectStructureError
		}

		id, err := s.QuestionRepository.CreateWithType(testID, type_, jsonData)

		if err != nil {
			log.Err(err).Send()
			return 0, make([]byte, 0), TestForbidden
		}

		return id, data.Ranges, nil
	default:
		return 0, make([]byte, 0), NotFountQuestionType
	}
}

func (s *QuestionService) Save(testID, userID, questionID int, type_ string, data []byte) error {
	_, err := s.TestRepository.Get(testID, userID)

	if err != nil {
		log.Err(err).Send()
		return TestForbidden
	}

	switch type_ {
	case Group:
		var group models.QuestionGroupData

		err = json.Unmarshal(data, &group)

		if err != nil {
			log.Err(err).Send()
			return InCorrectStructureError
		}

		err = s.QuestionRepository.Save(testID, questionID, type_, data)

		if err != nil {
			log.Err(err).Send()
			return InCorrectStructureError
		}

		return nil
	case Range:
		var range_ models.QuestionRangeData

		err = json.Unmarshal(data, &range_)

		if err != nil {
			log.Err(err).Send()
			return InCorrectStructureError
		}

		err = s.QuestionRepository.Save(testID, questionID, type_, data)

		if err != nil {
			log.Err(err).Send()
			return InCorrectStructureError
		}

		return nil
	default:
		return NotFountQuestionType
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
