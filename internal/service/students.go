package service

import (
	"App/internal/models"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
	"io"
)

var (
	ExcelOpenError            = errors.New("ошибка при открытии Excel файла")
	ExcelInCorrectFormatError = errors.New("некоррекный формат в таблице, проверьте структуру")

	StudentCreateManyError = errors.New("ошибка при сохранении студентов")
	StudentGetError        = errors.New("ошибка при получении студентов")
	StudentDeleteError     = errors.New("ошибка при удалении студента")
)

type StudentRepository interface {
	CreateMany(groupID int, students []models.Student) ([]models.Student, error)
	GetAll(groupID int) ([]models.Student, error)
	Delete(studentID, groupID int) error
}

type StudentService struct {
	StudentRepository
	GroupRepository
}

func NewStudentService(sr StudentRepository, rg GroupRepository) *StudentService {
	return &StudentService{sr, rg}
}

func (s *StudentService) Create(userID, groupID int, student models.Student) (int, error) {
	_, err := s.GroupRepository.Get(groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return 0, GroupNotFound
	}

	students, err := s.StudentRepository.CreateMany(groupID, []models.Student{student})

	if err != nil {
		log.Err(err).Send()
		return 0, StudentCreateManyError
	}

	return students[0].ID, nil
}

func (s *StudentService) CreateStudentFromExcel(userID, groupID int, fileContent io.Reader) ([]models.Student, error) {
	_, err := s.GroupRepository.Get(groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return nil, GroupNotFound
	}

	excel, err := excelize.OpenReader(fileContent)

	if err != nil {
		log.Err(err).Send()
		return nil, ExcelOpenError
	}

	rows, err := excel.GetRows(excel.GetSheetName(0))

	if err != nil {
		log.Err(err).Send()
		return nil, ExcelOpenError
	}

	var students []models.Student

	for _, row := range rows {
		var student models.Student

		countCells := len(row)

		switch countCells {
		case 3:
			student.Name = row[0]
			student.Surname = row[1]
			student.Patronymic = row[2]
		case 2:
			student.Name = row[0]
			student.Surname = row[1]
		default:
			return nil, ExcelInCorrectFormatError
		}

		students = append(students, student)
	}

	students, err = s.StudentRepository.CreateMany(groupID, students)

	if err != nil {
		log.Err(err).Send()
		return nil, StudentCreateManyError
	}

	return students, nil
}

func (s *StudentService) GetAll(userID, groupID int) ([]models.Student, error) {
	_, err := s.GroupRepository.Get(groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return nil, GroupNotFound
	}

	students, err := s.StudentRepository.GetAll(groupID)

	if err != nil {
		log.Err(err).Send()
		return nil, StudentGetError
	}

	return students, nil
}

func (s *StudentService) Delete(userID, groupID, studentID int) error {
	_, err := s.GroupRepository.Get(groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return GroupNotFound
	}

	err = s.StudentRepository.Delete(studentID, groupID)

	if err != nil {
		log.Err(err).Send()
		return StudentDeleteError
	}

	return nil
}
