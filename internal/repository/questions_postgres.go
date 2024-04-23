package repository

import (
	"App/internal/models"
	questions2 "App/internal/questions"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const DefaultTextQuestion = ""

const (
	Choose = "choose"
	Group  = "group"
	Range  = "range"

	CountGroups        = 3
	CountAnswerInGroup = 1
	CountRanges        = 3
)

var (
	NotDeleteRow = errors.New("ресурс не был удалён, хотя должен")
	NotUpdateRow = errors.New("ресурс не был обновлён, хотя должен")
	NotSaveError = errors.New("в слайсе нечего сохранять, len = 0")

	NotFoundTypeError = errors.New("неизвестный тип вопроса")
)

type QuestionPostgres struct {
	conn *sqlx.DB
}

func NewQuestionPostgres(conn *sqlx.DB) *QuestionPostgres {
	return &QuestionPostgres{conn}
}

func (r *QuestionPostgres) CreateWithType(testID int, type_ string, data []byte) (int, error) {
	text := models.GetTextFromType(type_)

	stmt := "INSERT INTO questions (text, test_id, data, type) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int

	err := r.conn.QueryRow(stmt, text, testID, data, type_).Scan(&id)

	if err != nil {
		log.Err(err).Send()
		return 0, err
	}

	return id, nil
}

func (r *QuestionPostgres) Save(testID, questionID int, type_ string, data []byte) error {
	stmt := "UPDATE questions SET data = $1 WHERE id = $2 AND test_id = $3 AND type = $4"

	result, err := r.conn.Exec(stmt, data, questionID, testID, type_)

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
	SELECT id, text, data, type
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
		var json string
		var type_ string

		if err := rows.Scan(&id, &text, &json, &type_); err != nil {
			log.Err(err).Send()
			continue
		}

		result := models.Question{
			ID:   id,
			Text: text,
			Data: json,
			Type: type_,
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

func (r *QuestionPostgres) GetAllWithAnswers(testID int) ([]models.QuestionWithAnswers, error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Info().Msg("возника паника при получении вопросов")
	//	}
	//}()

	query := `
	SELECT 
		q.id, 
		q.text, 
		q.type, 
		q.data, 
		a.id AS answer_id, 
		a.text AS answer_text, 
		a.is_correct AS is_answer_correct, 
		COALESCE(img.url, '') AS image_url
	FROM 
		questions q
	LEFT JOIN 
		answers a ON q.id = a.question_id
	LEFT JOIN (
		SELECT 
			id, 
			url, 
			question_id,
			ROW_NUMBER() OVER (PARTITION BY question_id ORDER BY id) AS rn
		FROM 
			images
	) AS img ON q.id = img.question_id AND img.rn = 1
	WHERE 
		q.test_id = $1
	ORDER BY 
		q.id ASC, 
		a.id ASC;
	`

	rows, err := r.conn.Query(query, testID)
	defer rows.Close()
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var questions []models.QuestionWithAnswers

	answersMap := make(map[int][]models.Answer)

	for rows.Next() {
		var questionID int
		var questionText string
		var questionType string
		var questionData []byte
		var questionImageURL string

		var answerID *int
		var answerText *string
		var isCorrect *bool

		if err := rows.Scan(
			&questionID, &questionText, &questionType, &questionData, &answerID, &answerText, &isCorrect, &questionImageURL,
		); err != nil {
			log.Err(err).Send()
			continue
		}

		if questionType != Choose {
			data, err := questions2.UnMarshalData(questionType, questionData)

			if err != nil {
				return nil, err
			}

			questions = append(questions, models.QuestionWithAnswers{
				ID:       questionID,
				Text:     questionText,
				Type:     questionType,
				Data:     data,
				ImageURL: questionImageURL,
				Answers:  []models.Answer{{ID: 0, Text: ""}},
			})

			continue
		}

		exists := false
		var q models.QuestionWithAnswers
		for _, item := range questions {
			if item.ID == questionID {
				q = item
				exists = true
				break
			}
		}

		if !exists {
			q = models.QuestionWithAnswers{
				ID:       questionID,
				Text:     questionText,
				Type:     questionType,
				ImageURL: questionImageURL,
			}
		}

		q.Answers = append(q.Answers, models.Answer{
			ID:        *answerID,
			Text:      *answerText,
			IsCorrect: *isCorrect,
		})

		if exists {
			for i, item := range questions {
				if item.ID == questionID {
					questions[i] = q
					break
				}
			}
		} else {
			questions = append(questions, q)
		}

		answersMap[questionID] = append(answersMap[questionID], models.Answer{
			ID:        *answerID,
			Text:      *answerText,
			IsCorrect: *isCorrect,
		})
	}

	return questions, nil
}

func (r *QuestionPostgres) CreateManyQuestions(
	testID int, questions []models.QuestionWithAnswersWithOutIsCorrect,
) ([]models.QuestionWithAnswersWithOutIsCorrect, error) {
	if len(questions) == 0 {
		return nil, NotSaveError
	}

	stmt := "INSERT INTO questions (text, test_id) VALUES "

	args := make([]any, len(questions)+1)

	args[0] = testID

	for i, question := range questions {
		stmt += fmt.Sprintf("($%d, $1)", i+2)

		args[i+1] = question.Text

		if i < len(questions)-1 {
			stmt += ", "
		}
	}

	stmt += " RETURNING id"

	var ids []int

	log.Info().Str("stmt", stmt).Send()
	log.Info().Any("args", args).Send()

	rows, err := r.conn.Query(stmt, args...)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Err(err).Send()
			return nil, err
		}
		ids = append(ids, id)
	}

	stmtA := "INSERT INTO answers (text, is_correct, question_id) VALUES "

	var argsA []any

	indexA := 1

	for i, question := range questions {
		for j, answer := range question.Answers {
			stmtA += fmt.Sprintf("($%d, $%d, $%d)", indexA, indexA+1, indexA+2)
			if j < len(question.Answers)-1 {
				stmtA += ", "
			}
			indexA += 3
			argsA = append(argsA, answer.Text, answer.IsCorrect, ids[i])
		}
		if i < len(questions)-1 {
			stmtA += ", "
		}
	}

	stmtA += " RETURNING id"

	log.Info().Str("stmtA", stmtA).Send()
	log.Info().Any("argsA", argsA).Send()

	rows, err = r.conn.Query(stmtA, argsA...)

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	var answerIDS []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Err(err).Send()
			return nil, err
		}

		answerIDS = append(answerIDS, id)
	}

	var index int

	for i, question := range questions {
		questions[i].ID = ids[i]
		for j, _ := range question.Answers {
			questions[i].Answers[j].ID = answerIDS[index]
			index++
		}
	}

	return questions, nil
}

func (r *QuestionPostgres) UploadImage(userID, testID, questionID int, filename string) error {
	tx, err := r.conn.Begin()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	stmt := "DELETE FROM images WHERE question_id = $1"

	_, err = tx.Exec(stmt, questionID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	stmt = "INSERT INTO images (url, question_id) VALUES ($1, $2)"

	_, err = tx.Exec(stmt, filename, questionID)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	err = tx.Commit()

	if err != nil {
		log.Err(err).Send()
		return err
	}

	return nil
}

func (r *QuestionPostgres) DeleteImage(questionID int) ([]string, error) {
	stmt := "DELETE FROM images WHERE question_id = $1 RETURNING url"

	var urls []string

	rows, err := r.conn.Query(stmt, questionID)

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return urls, nil
}
