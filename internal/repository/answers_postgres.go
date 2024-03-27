package repository

import (
	"App/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const DefaultTextAnswer = ""

type AnswerPostgres struct {
	conn *sqlx.DB
}

func NewAnswerPostgres(conn *sqlx.DB) *AnswerPostgres {
	return &AnswerPostgres{conn}
}

func (r *AnswerPostgres) Create(questionID int) (int, error) {
	stmt := "INSERT INTO answers (text, question_id) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.conn.QueryRow(stmt, DefaultTextAnswer, questionID).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}

func (r *AnswerPostgres) CreateThree(questionID int) ([]int, error) {
	stmt := `
	INSERT INTO answers (text, question_id)
	VALUES
		($1, $2),
		($3, $4),
		($5, $6)
	RETURNING id;
	`

	var ids []int

	rows, err := r.conn.Query(stmt, DefaultTextAnswer, questionID)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Err(err).Send()
			return nil, err
		}
		ids = append(ids, id)
	}

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return ids, nil
}

func (r *AnswerPostgres) GetAll(questionID int) ([]models.Answer, error) {
	query := `
	SELECT id, text, is_correct
	FROM answers
	WHERE question_id = $1
	ORDER BY id ASC
	`

	rows, err := r.conn.Query(query, questionID)

	defer rows.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var answers []models.Answer

	for rows.Next() {
		var id int
		var text string
		var isCorrect bool

		if err := rows.Scan(&id, &text, &isCorrect); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.Answer{
			ID:        id,
			Text:      text,
			IsCorrect: isCorrect,
		}
		answers = append(answers, result)
	}

	return answers, nil
}

func (r *AnswerPostgres) Update(answerUpdate models.AnswerUpdate, answerID, questionID int) error {
	stmt := `
	UPDATE answers
	SET text = $1, is_correct = $2
	WHERE id = $3 AND question_id = $4;
	`

	result, err := r.conn.Exec(
		stmt, answerUpdate.Text, answerUpdate.IsCorrect, answerID, questionID)

	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		return NotUpdateRow
	}

	return nil
}

func (r *AnswerPostgres) Delete(answerID, questionID int) error {
	stmt := `
	DELETE FROM answers WHERE id = $1 AND question_id = $2
	`

	result, err := r.conn.Exec(stmt, answerID, questionID)

	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		return NotDeleteRow
	}

	return nil
}
