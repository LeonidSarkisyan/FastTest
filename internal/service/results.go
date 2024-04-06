package service

import (
	"App/internal/models"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"slices"
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
	ResultRepository
}

func NewResultService(r ResultRepository) *ResultService {
	return &ResultService{r}
}

func (s *ResultService) GetResultByPassID(passID int) (models.ResultStudent, error) {
	result, err := s.ResultRepository.GetByPassID(passID)

	if err != nil {
		return models.ResultStudent{}, ResultCreateAlreadyError
	}

	return result, nil
}

func (s *ResultService) SaveResult(
	studentID, accessID, passID int, questions, questionsFromUser []models.QuestionWithAnswers,
	access models.AccessOut, timePass int,
) (models.ResultStudent, error) {

	log.Info().Any("questionsFromUser", questionsFromUser).Send()

	questionMap := make(map[int]models.QuestionForMap, len(questions))

	for _, q := range questions {
		answerMap := make(map[int]models.Answer, len(q.Answers))
		for _, a := range q.Answers {
			answerMap[a.ID] = a
		}
		questionMap[q.ID] = models.QuestionForMap{
			Question:   models.Question{ID: q.ID, Text: q.Text, Data: q.Data, Type: q.Type},
			AnswersMap: answerMap,
		}
	}

	var score int

	for _, qFromUser := range questionsFromUser {
		q, ok := questionMap[qFromUser.ID]

		if !ok {
			return models.ResultStudent{}, QuestionNotFoundError
		}

		var wasMistake bool

		switch qFromUser.Type {
		case Group:
			var data models.QuestionGroupData

			if err := mapstructure.Decode(qFromUser.Data, &data); err != nil {
				return models.ResultStudent{}, InCorrectStructureError
			}

			data.Groups = data.Groups[1:]

		Loop:
			for i, group := range data.Groups {
				groupFromDB := q.Data.(models.QuestionGroupData).Groups[i]

				if len(groupFromDB.Answers) != len(group.Answers) {
					wasMistake = true
					break Loop
				}

				for _, answer := range group.Answers {
					if !slices.Contains(groupFromDB.Answers, answer) {
						wasMistake = true
						break
					}
				}
			}
		default:
			for _, aFromUser := range qFromUser.Answers {
				a, ok := q.AnswersMap[aFromUser.ID]

				if !ok {
					return models.ResultStudent{}, AnswerNotFoundError
				}

				if a.IsCorrect != aFromUser.IsCorrect {
					wasMistake = true
					break
				}
			}
		}

		if !wasMistake {
			score += 1
		}
	}

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
		MaxScore: len(questions),
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
