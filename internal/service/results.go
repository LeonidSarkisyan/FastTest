package service

import (
	"App/internal/models"
	questions2 "App/internal/questions"
	"errors"
	"github.com/rs/zerolog/log"
	"strings"
)

var (
	QuestionNotFoundError = errors.New("при проверке не был обнаружен вопрос")

	AnswerNotFoundError = errors.New("при проверке не был обнаружен ответ")

	ResultCreateError        = errors.New("ошибка при сохранение результата")
	ResultGetError           = errors.New("ошибка при получении результата")
	ResultCreateAlreadyError = errors.New("вы уже завершали тест")
	ResultResetError         = errors.New("ошибка при обнулении результата")
)

type ResultRepository interface {
	Save(studentID, accessID, passID int, result models.ResultStudentIn) (int, error)
	Get(resultID int) (models.ResultStudent, error)
	GetAll(accessID int) ([]models.ResultStudent, error)
	ResetPass(passID int, access models.AccessOut) error

	GetByPassID(passID int) (models.ResultStudent, error)
}

type ResultService struct {
	*QuestionService
	ResultRepository
}

func NewResultService(qs *QuestionService, r ResultRepository) *ResultService {
	return &ResultService{qs, r}
}

func (s *ResultService) GetResultByPassID(passID int) (models.ResultStudent, error) {
	result, err := s.ResultRepository.GetByPassID(passID)

	if err != nil {
		return models.ResultStudent{}, ResultCreateAlreadyError
	}

	return result, nil
}

func (s *ResultService) SaveResult(
	studentID, accessID, passID int, questionsFromUser []models.QuestionWithAnswers,
	access models.AccessOut, timePass int,
) (models.ResultStudent, error) {

	questions := access.Questions.([]models.QuestionWithAnswers)
	score, err := questions2.Score(questions, questionsFromUser)

	if err != nil {
		log.Err(err).Send()
		return models.ResultStudent{}, err
	}

	maxScore := len(questionsFromUser)

	var mark int

	if score >= access.Criteria.Five {
		mark = 5
	} else if score >= access.Criteria.Four {
		mark = 4
	} else if score >= access.Criteria.Three {
		mark = 3
	} else {
		mark = 2
	}

	r := models.ResultStudentIn{
		Mark:     mark,
		Score:    score,
		MaxScore: maxScore,
		TimePass: timePass,
	}

	id, err := s.ResultRepository.Save(studentID, accessID, passID, r)

	if err != nil {
		if strings.HasPrefix(err.Error(), "pq: повторяющееся значение ключа нарушает ограничение уникальности") {
			return models.ResultStudent{}, ResultCreateAlreadyError
		}
		return models.ResultStudent{}, ResultCreateError
	}

	result, err := s.ResultRepository.Get(id)

	if err != nil {
		return models.ResultStudent{}, ResultGetError
	}

	return result, nil
}

func (s *ResultService) Reset(passID int, access models.AccessOut) error {
	err := s.ResultRepository.ResetPass(passID, access)

	if err != nil {
		return ResultResetError
	}

	return nil
}
