package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	questions2 "App/internal/questions"
	"App/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"math/rand/v2"
	"strconv"
	"time"
)

const (
	CheckBox    = "checkbox"
	RadioButton = "radio"
)

func shuffleAnswers(questions []models.QuestionWithAnswers) {
	for i := range questions {
		rand.Shuffle(len(questions[i].Answers), func(j, k int) {
			questions[i].Answers[j], questions[i].Answers[k] = questions[i].Answers[k], questions[i].Answers[j]
		})
	}
}

func (h *Handler) GetQuestionsForStudent(c *gin.Context) {
	studentIDCookie, err := c.Cookie("StudentID")

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	studentID, err := strconv.Atoi(studentIDCookie)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	accessID := MustID(c, "result_id")
	passID := MustID(c, "pass_id")

	pass, err := h.GetPassByStudentID(passID, studentID)

	if pass.IsActivated {
		SendErrorResponse(c, 401, "попытки кончились")
		return
	}

	access, err := h.TestService.GetResult(accessID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	startDate, err := time.Parse(time.RFC3339, access.DateStart)
	if err != nil {
		fmt.Println("Ошибка при разборе даты начала:", err)
		return
	}

	endDate, err := time.Parse(time.RFC3339, access.DateEnd)
	if err != nil {
		fmt.Println("Ошибка при разборе даты конца:", err)
		return
	}

	if err = utils.CheckDateLimit(startDate, endDate); err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	_, err = h.TestService.Get(access.TestID, access.UserID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	log.Info().Any("question", access.Questions).Send()

	questions := access.Questions.([]models.QuestionWithAnswers)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	if access.Shuffle {
		dest := make([]models.QuestionWithAnswers, len(questions))
		perm := rand.Perm(len(questions))
		for i, v := range perm {
			dest[v] = questions[i]
		}
		questions = dest
	}

	questions2.HideData(questions)

	err = h.TestService.ClosePass(passID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	shuffleAnswers(questions)

	c.JSON(200, gin.H{
		"test_id":   access.TestID,
		"access":    access,
		"questions": responses.NewListResponse(questions),
	})
}

func (h *Handler) CreateResult(c *gin.Context) {
	studentIDCookie, err := c.Cookie("StudentID")

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	studentID, err := strconv.Atoi(studentIDCookie)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	accessID := MustID(c, "result_id")
	passID := MustID(c, "pass_id")

	_, err = h.GetPassByStudentID(passID, studentID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	access, err := h.TestService.GetResult(accessID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	_, err = h.TestService.Get(access.TestID, access.UserID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	questions := access.Questions

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	var r models.Result

	if err := c.Bind(&r); err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	var questionWithAnswer []models.QuestionWithAnswers

	for _, q := range r.Questions {
		newQ := models.QuestionWithAnswers{
			ID:   q.ID,
			Text: q.Text,
			Data: q.Data,
			Type: q.Type,
		}

		for _, a := range q.Answers {
			newQ.Answers = append(newQ.Answers, models.Answer{
				ID:        a.ID,
				Text:      a.Text,
				IsCorrect: a.IsCorrect,
			})
		}

		questionWithAnswer = append(questionWithAnswer, newQ)
	}

	log.Info().Any("questions тут!", questions).Send()

	result, err := h.ResultService.SaveResult(
		studentID, accessID, passID, questionWithAnswer, access, r.TimePass,
	)

	if err != nil {
		SendErrorResponse(c, 409, err.Error())
		return
	}

	ch, ok := h.ClientManager.PassesMap.Load(passID)

	if ok {
		ch.(chan models.ResultStudent) <- result
	}

	resultCh, ok := h.ClientManager.ResultsMap.Load(accessID)

	if ok {
		resultCh.(chan models.ResultStudent) <- result
	}

	c.JSON(201, gin.H{
		"result": result,
	})
}

func (h *Handler) AbortPage(c *gin.Context) {
	c.HTML(409, "error.html", gin.H{
		"error": "Ошибка 409, прохождение теста было прервано создателем теста",
	})
}
