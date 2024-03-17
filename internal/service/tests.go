package service

import (
	"App/internal/models"
	"App/pkg/utils"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

var (
	TestCreateError = errors.New("ошибка при создании теста")
	TestGetError    = errors.New("ошибка при получении теста")
	TestUpdateError = errors.New("ошибка при обновлении теста")

	TestAccessCreateError = errors.New("ошибка при создании доступа к тесту")

	InCorrectCriteria   = errors.New("некорректная структура критерия")
	GeneratePassesError = errors.New("ошибка при генерации пропусков")

	NotStudentError = errors.New("вы выбрали группу, где нет студентов")

	AccessGetError = errors.New("ошибка при получении результатов теста")
)

type TestRepository interface {
	Create(test models.Test) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(testID, userID int, title string) error

	CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error)
	GetAccess(userID, accessID int) (models.AccessOut, error)

	CreateManyPasses(accessID int, passes []models.PassesIn) error
}

type TestService struct {
	TestRepository
	StudentRepository
	*QuestionService
	*GroupService
}

func NewTestService(r TestRepository, sr StudentRepository, sq *QuestionService, gs *GroupService) *TestService {
	return &TestService{r, sr, sq, gs}
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

	log.Info().Int("test_id", id).Send()

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

func (s *TestService) CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error) {
	_, err := s.Get(testID, userID)

	if err != nil {
		return 0, err
	}

	_, err = s.GroupService.Get(groupID, userID)

	if err != nil {
		return 0, err
	}

	criteriaJson, err := json.Marshal(accessIn.Criteria)

	if err != nil {
		return 0, InCorrectCriteria
	}

	accessIn.CriteriaJson = criteriaJson

	students, err := s.StudentRepository.GetAll(groupID)

	if err != nil {
		return 0, StudentGetError
	}

	if len(students) == 0 {
		return 0, NotStudentError
	}

	err = s.CheckTest(testID)

	if err != nil {
		return 0, err
	}

	id, err := s.TestRepository.CreateAccess(userID, testID, groupID, accessIn)

	if err != nil {
		log.Err(err).Send()
		return 0, TestAccessCreateError
	}

	err = s.CreatePasses(groupID, id)

	if err != nil {
		log.Err(err).Send()
		return 0, GeneratePassesError
	}

	return id, nil
}

func (s *TestService) CreatePasses(groupID, accessID int) error {
	students, err := s.StudentRepository.GetAll(groupID)

	if err != nil {
		log.Err(err).Send()
		return StudentGetError
	}

	countStudents := len(students)

	passes := make([]models.PassesIn, countStudents)

	codes := utils.GenerateSixDigitNumber(countStudents)

	for i, student := range students {
		passes[i].Code = codes[i]
		passes[i].StudentID = student.ID
	}

	err = s.TestRepository.CreateManyPasses(accessID, passes)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	return nil
}

func (s *TestService) CheckTest(testID int) error {
	questions, err := s.QuestionRepository.GetAllWithAnswers(testID)

	if err != nil {
		return QuestionGetCreate
	}

	for i, q := range questions {
		n := i + 1

		if len(q.Text) == 0 {
			return fmt.Errorf("у вопроса с номером %d нет текста, проверьте тест", n)
		}

		if len(q.Answers) < 2 {
			return fmt.Errorf("у вопроса с номером %d меньше двух вариантов ответа, проверьте тест", n)
		}

		var isCorrect bool

		for _, a := range q.Answers {
			if len(a.Text) == 0 {
				return fmt.Errorf("у варианта ответа под вопросом с номером %d нет текста, проверьте тест", n)
			}

			if a.IsCorrect {
				isCorrect = a.IsCorrect
			}
		}

		if !isCorrect {
			return fmt.Errorf("у вопроса с номером %d нет хотя бы одного правильного ответа, проверьте тест", n)
		}
	}

	return nil
}

func (s *TestService) GetAccess(userID, accessID int) (models.AccessOut, error) {
	a, err := s.TestRepository.GetAccess(userID, accessID)

	if err != nil {
		log.Err(err).Send()
		return models.AccessOut{}, AccessGetError
	}

	return a, nil
}
