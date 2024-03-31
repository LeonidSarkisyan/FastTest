package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
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

	test, err := h.TestService.Get(access.TestID, access.UserID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	questions, err := h.QuestionService.GetAllQuestionsWithAnswers(test.ID)

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

	for i, q := range questions {
		var countRight int
		for _, a := range q.Answers {
			if a.IsCorrect {
				countRight++
			}
		}
		if countRight >= 2 {
			questions[i].Type = CheckBox
		} else {
			questions[i].Type = RadioButton
		}
	}

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

	resetChannel := make(chan int)
	timesChannel := make(chan int)

	h.ClientManager.ResetMap.Store(passID, resetChannel)
	h.ClientManager.TimesMap.Store(passID, timesChannel)

	go func() {
		timeout := time.Duration(access.PassageTime)*time.Minute + 5*time.Second
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()

		for {
			resetCh, _ := h.ClientManager.ResetMap.Load(passID)

			select {
			case <-ticker.C:
				resultStudent, err := h.ResultService.SaveResult(
					studentID, accessID, passID, questions, []models.QuestionWithAnswers{}, access, access.PassageTime*60,
				)
				if err != nil {
					log.Err(err).Send()
					return
				}
				h.SendToBroadcast(Message{
					UserID: access.UserID,
					Result: resultStudent,
				})

				resetCh, ok := h.ClientManager.ResetMap.Load(passID)
				if ok {
					resetCh.(chan int) <- 1
				}
				timesCh, ok := h.ClientManager.TimesMap.Load(passID)
				if ok {
					timesCh.(chan int) <- 1
				}
				return
			case <-resetCh.(chan int):
				log.Info().Msg("тест прерван или завершён, выключаем принудельную двойку")
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		msg := Message{
			UserID: access.UserID,
			PassID: passID,
			Result: models.ResultStudent{
				Mark:         -1,
				DateTimePass: time.Time{},
				PassID:       passID,
				TimePass:     1,
			},
		}

		for {
			timesCh, _ := h.ClientManager.TimesMap.Load(passID)

			select {
			case <-ticker.C:
				log.Info().Msg("прошла секунда, отправляю...")
				h.ClientManager.SendToBroadcast(msg)
				msg.PassID = 0
				h.ClientManager.SendToBroadcast(msg)
				msg.Result.TimePass++
				msg.PassID = passID
			case <-timesCh.(chan int):
				log.Info().Msg("тест завершён, выключает посекундное обновление")
				return
			}
		}
	}()
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

	test, err := h.TestService.Get(access.TestID, access.UserID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	questions, err := h.QuestionService.GetAllQuestionsWithAnswers(test.ID)

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

	result, err := h.ResultService.SaveResult(
		studentID, accessID, passID, questions, questionWithAnswer, access, r.TimePass,
	)

	if err != nil {
		SendErrorResponse(c, 409, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"result": result,
	})

	go func() {
		message := Message{
			UserID: access.UserID,
			Result: result,
		}

		h.ClientManager.SendToBroadcast(message)

		if val, ok := h.ClientManager.TimesMap.Load(passID); ok {
			if timesChannel, ok := val.(chan int); ok {
				timesChannel <- 1
			}
		}

		if val, ok := h.ClientManager.ResetMap.Load(passID); ok {
			if resetChannel, ok := val.(chan int); ok {
				resetChannel <- 1
			}
		}
	}()

}

func (h *Handler) AbortPage(c *gin.Context) {
	c.HTML(409, "error.html", gin.H{
		"error": "Ошибка 409, прохождение теста было прервано создателем теста",
	})
}
