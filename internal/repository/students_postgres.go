package repository

import (
	"App/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type StudentPostgres struct {
	db *sqlx.DB
}

func NewStudentPostgres(db *sqlx.DB) *StudentPostgres {
	return &StudentPostgres{db}
}

func (r *StudentPostgres) CreateMany(groupID int, students []models.Student) ([]models.Student, error) {
	stmt := "INSERT INTO students (name, surname, patronymic, group_id) VALUES "

	for i, student := range students {
		stmt += fmt.Sprintf("('%s', '%s', '%s', %d)", student.Name, student.Surname, student.Patronymic, groupID)

		if i != len(students)-1 {
			stmt += ","
		}
	}

	stmt += " RETURNING id;"

	rows, err := r.db.Query(stmt)

	index := 0
	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		students[index].ID = id
		index++
	}

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return students, nil
}

func (r *StudentPostgres) GetAll(groupID int) ([]models.Student, error) {
	query := `
	SELECT id, name, surname, patronymic 
	FROM students 
	WHERE group_id = $1
	ORDER BY surname
	`

	var students []models.Student

	rows, err := r.db.Query(query, groupID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var student models.Student

		if err := rows.Scan(&student.ID, &student.Name, &student.Surname, &student.Patronymic); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (r *StudentPostgres) Delete(studentID, groupID int) error {
	stmt := `
	DELETE FROM students WHERE id = $1 AND group_id = $2
	`

	result, err := r.db.Exec(stmt, studentID, groupID)

	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		return NotDeleteRow
	}

	return nil
}

func (r *StudentPostgres) Get(studentID int) (models.Student, error) {
	query := `
	SELECT id, name, surname, patronymic 
	FROM students 
	WHERE id = $1
	`

	var s models.Student

	err := r.db.QueryRow(query, studentID).Scan(&s.ID, &s.Name, &s.Surname, &s.Patronymic)

	if err != nil {
		log.Err(err).Send()
		return models.Student{}, err
	}

	return s, nil
}
