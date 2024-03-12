package repository

import (
	"App/internal/models"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const DefaultTextQuestion = ""

var (
	NotDeleteRow = errors.New("ресурс не был удалён, хотя должен")
	NotUpdateRow = errors.New("ресурс не был обновлён, хотя должен")
)

type QuestionPostgres struct {
	conn *sqlx.DB
}

func NewQuestionPostgres(conn *sqlx.DB) *QuestionPostgres {
	return &QuestionPostgres{conn}
}

func (r *QuestionPostgres) Create(testID int) (int, error) {
	stmt := "INSERT INTO questions (text, test_id) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.conn.QueryRow(stmt, DefaultTextQuestion, testID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *QuestionPostgres) GetAll(testID int) ([]models.Question, error) {
	query := `
	SELECT id, text
	FROM questions
	WHERE test_id = $1
	ORDER BY id ASC
	`

	rows, err := r.conn.Query(query, testID)

	defer rows.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var questions []models.Question

	for rows.Next() {
		var id int
		var text string

		if err := rows.Scan(&id, &text); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.Question{
			ID:   id,
			Text: text,
		}
		questions = append(questions, result)
	}

	return questions, nil
}

func (r *QuestionPostgres) Update(questionID, testID int, question models.QuestionUpdate) error {
	stmt := `
	UPDATE questions
	SET text = $1
	WHERE id = $2 AND test_id = $3;
	`

	result, err := r.conn.Exec(stmt, question.Text, questionID, testID)

	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		return NotUpdateRow
	}

	return nil
}

func (r *QuestionPostgres) Get(questionID, testID int) (models.Question, error) {
	query := `
	SELECT id, text
	FROM questions
	WHERE test_id = $1 AND id = $2;
	`

	var question models.Question

	err := r.conn.QueryRow(query, testID, questionID).Scan(&question.ID, &question.Text)

	if err != nil {
		log.Err(err).Send()
		return models.Question{}, err
	}

	return question, nil
}

func (r *QuestionPostgres) Delete(questionID, testID int) error {
	stmt := `
	DELETE FROM questions WHERE id = $1 AND test_id = $2
	`

	result, err := r.conn.Exec(stmt, questionID, testID)

	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		return NotDeleteRow
	}

	return nil
}
