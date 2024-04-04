package repository

import (
	"App/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type ResultPostgres struct {
	db *sqlx.DB
}

func NewResultPostgres(db *sqlx.DB) *ResultPostgres {
	return &ResultPostgres{db}
}

func (r *ResultPostgres) Save(studentID, accessID, passID int, result models.ResultStudentIn) (int, error) {
	stmt := `
	INSERT INTO results (mark, score, max_score, pass_id, access_id, student_id, time_pass)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;
	`

	var id int

	log.Info().Int("access_id", accessID).Int("studentID", studentID).Int("passID", passID).Send()

	err := r.db.QueryRow(
		stmt, result.Mark, result.Score, result.MaxScore, passID, accessID, studentID, result.TimePass,
	).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}

func (r *ResultPostgres) GetByPassID(passID int) (models.ResultStudent, error) {
	stmt := `
	SELECT mark, score, max_score, pass_id, access_id, student_id, time_pass
	FROM results
	WHERE pass_id = $1;
	`

	var result models.ResultStudent

	row := r.db.QueryRow(stmt, passID)
	err := row.Scan(
		&result.Mark, &result.Score, &result.MaxScore, &result.PassID, &result.AccessID, &result.StudentID,
		&result.TimePass,
	)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Err(err).Send()
		return result, err
	}

	return result, nil
}

func (r *ResultPostgres) Get(resultID int) (models.ResultStudent, error) {
	stmt := `
	SELECT mark, score, max_score, pass_id, access_id, student_id, time_pass
	FROM results
	WHERE id = $1;
	`

	var result models.ResultStudent

	row := r.db.QueryRow(stmt, resultID)
	err := row.Scan(
		&result.Mark, &result.Score, &result.MaxScore, &result.PassID, &result.AccessID, &result.StudentID,
		&result.TimePass,
	)
	if err != nil {
		log.Err(err).Send()
		return models.ResultStudent{}, err
	}

	return result, nil
}

func (r *ResultPostgres) GetAll(accessID int) ([]models.ResultStudent, error) {
	stmt := `
	SELECT mark, score, max_score, pass_id, access_id, student_id, time_pass
	FROM results
	WHERE access_id = $1;
	`

	var results []models.ResultStudent

	row, err := r.db.Query(stmt, accessID)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for row.Next() {
		var res models.ResultStudent

		err = row.Scan(
			&res.Mark, &res.Score, &res.MaxScore, &res.PassID, &res.AccessID, &res.StudentID, &res.TimePass,
		)

		if err != nil {
			log.Err(err).Send()
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
}

func (r *ResultPostgres) ResetPass(passID int, access models.AccessOut) error {
	tx, err := r.db.Begin()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	stmt := "DELETE FROM results WHERE pass_id = $1;"

	result, err := tx.Exec(stmt, passID)

	_, err = result.RowsAffected()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	stmt = "UPDATE passes SET is_activated = false, datetime_activate = null WHERE id = $1;"

	result, err = tx.Exec(stmt, passID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	if count == 0 {
		return NotUpdateRow
	}

	err = tx.Commit()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	return nil
}
