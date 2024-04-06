package repository

import (
	"App/internal/models"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"strings"
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
	SELECT t.id, t.title, COALESCE(COUNT(q.id), 0), is_deleted
	FROM tests t 
	LEFT JOIN questions q ON t.id = q.test_id
	WHERE t.id = $1 AND t.user_id = $2
	GROUP BY t.id, t.title
	`

	var testOut models.TestOut

	err := r.conn.QueryRow(query, testID, userID).Scan(&testOut.ID, &testOut.Title, &testOut.Count, &testOut.IsDeleted)

	if err != nil {
		log.Err(err).Send()
		return models.TestOut{}, err
	}

	return testOut, nil
}

func (r *TestPostgres) GetIfNotDelete(testID, userID int) (models.TestOut, error) {
	query := `
	SELECT t.id, t.title, COALESCE(COUNT(q.id), 0), is_deleted
	FROM tests t 
	LEFT JOIN questions q ON t.id = q.test_id
	WHERE t.id = $1 AND t.user_id = $2 AND is_deleted = false
	GROUP BY t.id, t.title
	`

	var testOut models.TestOut

	err := r.conn.QueryRow(query, testID, userID).Scan(&testOut.ID, &testOut.Title, &testOut.Count, &testOut.IsDeleted)

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
	WHERE user_id = $1 AND is_deleted = false
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

func (r *TestPostgres) Delete(userID, testID int) error {
	query := "UPDATE tests SET is_deleted = true WHERE user_id = $1 AND id = $2"

	res, err := r.conn.Exec(query, userID, testID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	if count == 0 {
		return NotDeleteRow
	}

	return nil
}

func (r *TestPostgres) CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error) {
	stmt := `
	INSERT INTO accesses (shuffle, date_start, date_end, passage_time, criteria, user_id, test_id, group_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
	`

	var id int

	err := r.conn.QueryRow(
		stmt, accessIn.Shuffle, accessIn.DateStart, accessIn.DateEnd,
		accessIn.PassageTime, accessIn.CriteriaJson, userID, testID, groupID,
	).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}

func (r *TestPostgres) CreateManyPasses(accessID int, passes []models.PassesIn) error {
	stmt := "INSERT INTO passes (code, access_id, student_id) VALUES "

	for i, p := range passes {
		stmt += fmt.Sprintf("(%d, %d, %d)", p.Code, accessID, p.StudentID)

		if i != len(passes)-1 {
			stmt += ","
		}
	}

	_, err := r.conn.Exec(stmt)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	return nil
}

func (r *TestPostgres) GetAccess(userID, accessID int) (models.AccessOut, error) {
	query := `
	SELECT id, test_id, group_id, date_start, date_end, passage_time 
	FROM accesses
	WHERE id = $1 AND user_id = $2;
	`

	var a models.AccessOut

	err := r.conn.QueryRow(query, accessID, userID).Scan(
		&a.ID, &a.TestID, &a.GroupID, &a.DateStart, &a.DateEnd, &a.PassageTime,
	)

	if err != nil {
		log.Err(err).Send()
		return models.AccessOut{}, err
	}

	return a, nil
}

func (r *TestPostgres) GetPasses(resultID int) ([]models.Passes, error) {
	query := `
	SELECT passes.id, passes.is_activated, passes.datetime_activate, passes.code, passes.student_id
	FROM passes
	JOIN students ON passes.student_id = students.id
	WHERE passes.access_id = $1
	ORDER BY students.surname;
	`

	var passes []models.Passes

	rows, err := r.conn.Query(query, resultID)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for rows.Next() {
		var p models.Passes

		if err := rows.Scan(&p.ID, &p.IsActivated, &p.DateTimeActivated, &p.Code, &p.StudentID); err != nil {
			log.Err(err).Send()
			continue
		}
		passes = append(passes, p)
	}

	return passes, nil
}

func (r *TestPostgres) GetAllAccesses(userID int) ([]models.AccessOut, error) {
	query := `
	SELECT a.id, a.test_id, a.group_id, a.date_start, a.date_end, a.passage_time, t.title, g.name 
	FROM accesses a
	JOIN tests AS t ON a.test_id = t.id 
	JOIN groups AS g ON a.group_id = g.id 
	WHERE a.user_id = $1
	ORDER BY a.id DESC;
	`

	var accesses []models.AccessOut

	rows, err := r.conn.Query(query, userID)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for rows.Next() {
		var a models.AccessOut

		if err := rows.Scan(
			&a.ID, &a.TestID, &a.GroupID, &a.DateStart, &a.DateEnd, &a.PassageTime, &a.Test.Title,
			&a.GroupOut.Name,
		); err != nil {
			log.Err(err).Send()
			return nil, err
		}

		a.DateStart = strings.ReplaceAll(
			strings.ReplaceAll(a.DateStart, "T00:00:00Z", ""), "-", ".")

		a.DateEnd = strings.ReplaceAll(
			strings.ReplaceAll(a.DateEnd, "T00:00:00Z", ""), "-", ".")

		accesses = append(accesses, a)
	}

	return accesses, nil
}

func (r *TestPostgres) GetPass(resultID int, code int64) (models.Passes, error) {
	query := `
	SELECT id, is_activated, datetime_activate, code, student_id
	FROM passes
	WHERE access_id = $1 AND code = $2
	`

	var p models.Passes

	err := r.conn.QueryRow(query, resultID, code).Scan(
		&p.ID, &p.IsActivated, &p.DateTimeActivated, &p.Code, &p.StudentID,
	)

	if err != nil {
		log.Err(err).Send()
		return models.Passes{}, err
	}

	return p, nil
}

func (r *TestPostgres) GetResult(resultID int) (models.AccessOut, error) {
	query := `
	SELECT id, test_id, group_id, date_start, date_end, passage_time, shuffle, user_id, criteria
	FROM accesses
	WHERE id = $1;
	`

	var a models.AccessOut

	err := r.conn.QueryRow(query, resultID).Scan(
		&a.ID, &a.TestID, &a.GroupID, &a.DateStart, &a.DateEnd, &a.PassageTime, &a.Shuffle, &a.UserID, &a.CriteriaJson,
	)

	if err != nil {
		log.Err(err).Send()
		return models.AccessOut{}, err
	}

	err = json.Unmarshal(a.CriteriaJson, &a.Criteria)

	if err != nil {
		log.Err(err).Send()
		return models.AccessOut{}, err
	}

	return a, nil
}

func (r *TestPostgres) GetPassByStudentID(passID, studentID int) (models.Passes, error) {
	query := `
	SELECT id, is_activated, datetime_activate, code, student_id
	FROM passes
	WHERE id = $1 AND student_id = $2
	`

	var p models.Passes

	err := r.conn.QueryRow(query, passID, studentID).Scan(
		&p.ID, &p.IsActivated, &p.DateTimeActivated, &p.Code, &p.StudentID,
	)

	if err != nil {
		log.Err(err).Send()
		return models.Passes{}, err
	}

	return p, nil
}

func (r *TestPostgres) ClosePass(passID int) error {
	stmt := "UPDATE passes SET is_activated = true, datetime_activate = CURRENT_TIMESTAMP WHERE id = $1;"

	result, err := r.conn.Exec(stmt, passID)

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

	return nil
}
