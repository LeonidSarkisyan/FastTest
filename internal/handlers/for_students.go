package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"math/rand/v2"
	"strconv"
	"sync"
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

	go func() {
		var mu sync.Mutex

		mu.Lock()
		h.ClientManager.ResetMap[passID] = make(chan int)
		mu.Unlock()

		timeout := time.Duration(access.PassageTime)*time.Minute + 5*time.Second

		for {
			select {
			case <-time.After(timeout):
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

				go func() {
					var mu sync.Mutex
					mu.Lock()
					h.ClientManager.TimesMap[passID] <- 1
					h.ClientManager.ResetMap[passID] <- 1
					mu.Unlock()
				}()
			case <-h.ClientManager.ResetMap[passID]:
				var mu sync.Mutex
				log.Info().Msg("тест прерван или завершён, выключаем принудельную двойку")
				mu.Lock()
				delete(h.ClientManager.TimesMap, passID)
				mu.Unlock()
				return
			}
		}

	}()

	go func() {
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

		var mu sync.Mutex

		mu.Lock()
		h.ClientManager.TimesMap[passID] = make(chan int)
		mu.Unlock()

		for {
			select {
			case <-time.After(time.Second):
				log.Info().Msg("прошла секунда, отправляю...")
				h.ClientManager.SendToBroadcast(msg)
				msg.PassID = 0
				h.ClientManager.SendToBroadcast(msg)
				msg.Result.TimePass++
				msg.PassID = passID
			case <-h.ClientManager.TimesMap[passID]:
				log.Info().Msg("тест завершён, выключает посекундное обновление")
				mu.Lock()
				delete(h.ClientManager.TimesMap, passID)
				mu.Unlock()
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

		var mu sync.Mutex
		mu.Lock()
		h.ClientManager.TimesMap[passID] <- 1
		h.ClientManager.ResetMap[passID] <- 1
		mu.Unlock()
	}()
}

func (h *Handler) AbortPage(c *gin.Context) {
	c.HTML(409, "error.html", gin.H{
		"error": "Ошибка 409, прохождение теста было прервано создателем теста",
	})
}
