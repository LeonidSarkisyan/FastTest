package repository

import (
	"App/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type TestPostgres struct {
	conn *sqlx.DB
}

func NewTestPostgres(conn *sqlx.DB) *TestPostgres {
	return &TestPostgres{conn}
}

func (r *TestPostgres) Create(test models.Test) (int, error) {
	stmt := "INSERT INTO tests (title, user_id) VALUES ($1, $2) RETURNING id"

	var id int

	err := r.conn.QueryRow(stmt, test.Title, test.UserID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TestPostgres) Get(testID, userID int) (models.TestOut, error) {
	query := "SELECT id, title FROM tests WHERE id = $1 AND user_id = $2"

	var testOut models.TestOut

	err := r.conn.QueryRow(query, testID, userID).Scan(&testOut.ID, &testOut.Title)

	if err != nil {
		log.Err(err).Send()
		return models.TestOut{}, err
	}

	return testOut, nil
}

func (r *TestPostgres) GetAll(userID int) ([]models.TestOut, error) {
	query := `
	SELECT id, title, EXTRACT(EPOCH FROM datetime_create)::BIGINT
	FROM tests
	WHERE user_id = $1
	ORDER BY datetime_create DESC;
	`

	rows, err := r.conn.Query(query, userID)

	defer rows.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var tests []models.TestOut

	for rows.Next() {
		var id int
		var title string
		var dateTimeCreate int64

		if err := rows.Scan(&id, &title, &dateTimeCreate); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.TestOut{
			ID:             id,
			Title:          title,
			DateTimeCreate: dateTimeCreate * 1000,
		}
		tests = append(tests, result)
	}

	return tests, nil
}

func (r *TestPostgres) UpdateTitle(testID, userID int, title string) error {
	query := "UPDATE tests SET title = $1 WHERE id = $2 AND user_id = $3"

	result, err := r.conn.Exec(query, title, testID, userID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	if count, err := result.RowsAffected(); count == 0 || err != nil {
		log.Error().Msg("Rows affected: " + string(count))
		log.Err(err).Send()
		return err
	}

	return nil
}
