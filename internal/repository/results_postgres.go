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
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`

	var id int

	log.Info().Int("access_id", accessID).Int("studentID", studentID).Int("passID", passID).Send()

	err := r.db.QueryRow(stmt, result.Mark, result.Score, result.MaxScore, passID, accessID, studentID).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}

func (r *ResultPostgres) Get(resultID int) (models.ResultStudent, error) {
	stmt := `
	SELECT mark, score, max_score, pass_id, access_id, student_id
	FROM results
	WHERE id = $1;
	`

	var result models.ResultStudent

	row := r.db.QueryRow(stmt, resultID)
	err := row.Scan(&result.Mark, &result.Score, &result.MaxScore, &result.PassID, &result.AccessID, &result.StudentID)
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
