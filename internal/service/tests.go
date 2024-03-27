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
	PassesGetError = errors.New("ошибка при получении пропусков")

	PassGetError  = errors.New("ошибка при получении пропуска")
	PassNotFound  = errors.New("неверный код")
	PassDontClose = errors.New("не удалось закрыть доступ к тесту")
)

type TestRepository interface {
	Create(test models.Test) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(testID, userID int, title string) error

	CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error)
	GetAccess(userID, accessID int) (models.AccessOut, error)
	GetAllAccesses(userID int) ([]models.AccessOut, error)
	GetResult(resultID int) (models.AccessOut, error)

	CreateManyPasses(accessID int, passes []models.PassesIn) error
	GetPasses(resultID int) ([]models.Passes, error)
	GetPass(resultID int, code int64) (models.Passes, error)
	GetPassByStudentID(passID, studentID int) (models.Passes, error)
	ClosePass(passID int) error
}

type TestService struct {
	TestRepository
	StudentRepository
	*QuestionService
	*GroupService
	ResultRepository
}

func NewTestService(
	r TestRepository, sr StudentRepository, sq *QuestionService, gs *GroupService, rr ResultRepository,
) *TestService {
	return &TestService{r, sr, sq, gs, rr}
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

	_, _, err = s.QuestionService.Create(id, userID)

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

func (s *TestService) GetPassesAndStudents(resultID, userID int) (models.ForResultTable, error) {
	access, err := s.TestRepository.GetAccess(userID, resultID)

	if err != nil {
		log.Err(err).Send()
		return models.ForResultTable{}, AccessGetError
	}

	_, err = s.GroupService.Get(access.GroupID, userID)

	if err != nil {
		log.Err(err).Send()
		return models.ForResultTable{}, GroupGetError
	}

	passes, err := s.TestRepository.GetPasses(resultID)

	if err != nil {
		log.Err(err).Send()
		return models.ForResultTable{}, PassesGetError
	}

	students, err := s.StudentRepository.GetAll(access.GroupID)

	if err != nil {
		log.Err(err).Send()
		return models.ForResultTable{}, StudentGetError
	}

	forResultTable := models.ForResultTable{
		Students: students,
		Passes:   passes,
		Results:  make([]models.ResultStudent, len(passes)),
	}

	results, err := s.ResultRepository.GetAll(access.ID)

	if err != nil {
		log.Err(err).Send()
		return models.ForResultTable{}, ResultGetError
	}

	for i, p := range forResultTable.Passes {
		for _, r := range results {
			if r.PassID == p.ID {
				forResultTable.Results[i] = r
			}
		}
	}

	return forResultTable, nil
}

func (s *TestService) GetAllAccessess(userID int) ([]models.AccessOut, error) {
	accesses, err := s.TestRepository.GetAllAccesses(userID)

	if err != nil {
		log.Err(err).Send()
		return nil, AccessGetError
	}

	return accesses, nil
}

func (s *TestService) GetPassByCode(resultID int, code int64) (models.Passes, error) {
	pass, err := s.TestRepository.GetPass(resultID, code)

	if err != nil {
		log.Err(err).Send()

		if err.Error() == "sql: no rows in result set" {
			return models.Passes{}, PassNotFound
		}

		return models.Passes{}, PassGetError
	}

	return pass, nil
}

func (s *TestService) GetResult(resultID int) (models.AccessOut, error) {
	r, err := s.TestRepository.GetResult(resultID)

	if err != nil {
		return models.AccessOut{}, AccessGetError
	}

	return r, nil
}

func (s *TestService) GetPassByStudentID(passID, studentID int) (models.Passes, error) {
	pass, err := s.TestRepository.GetPassByStudentID(passID, studentID)

	if err != nil {
		return models.Passes{}, PassGetError
	}

	return pass, nil
}

func (s *TestService) ClosePass(passID int) error {
	err := s.TestRepository.ClosePass(passID)

	if err != nil {
		return PassDontClose
	}

	return nil
}
