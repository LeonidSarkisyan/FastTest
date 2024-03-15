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
	query := `
	SELECT t.id, t.title, COALESCE(COUNT(q.id), 0)
	FROM tests t 
	LEFT JOIN questions q ON t.id = q.test_id
	WHERE t.id = $1 AND t.user_id = $2
	GROUP BY t.id, t.title
	`

	var testOut models.TestOut

	err := r.conn.QueryRow(query, testID, userID).Scan(&testOut.ID, &testOut.Title, &testOut.Count)

	if err != nil {
		log.Err(err).Send()
		return models.TestOut{}, err
	}

	return testOut, nil
}

func (r *TestPostgres) GetAll(userID int) ([]models.TestOut, error) {
	query := `
	SELECT t.id, t.title, EXTRACT(EPOCH FROM t.datetime_create)::BIGINT, COALESCE(COUNT(questions.id), 0)
	FROM tests t
	LEFT JOIN questions ON t.id = questions.test_id
	WHERE user_id = $1
	GROUP BY t.id, t.title
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
		var count int

		if err := rows.Scan(&id, &title, &dateTimeCreate, &count); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.TestOut{
			ID:             id,
			Title:          title,
			DateTimeCreate: dateTimeCreate * 1000,
			Count:          count,
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

func (r *TestPostgres) CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error) {
	stmt := `
	INSERT INTO accesses (date_start, date_end, passage_time, criteria, user_id, test_id, group_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;
	`

	var id int

	err := r.conn.QueryRow(
		stmt, accessIn.DateStart, accessIn.DateEnd, accessIn.PassageTime, accessIn.CriteriaJson, userID, testID, groupID,
	).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}
