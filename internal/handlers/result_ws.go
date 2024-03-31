package handlers

import (
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (h *Handler) CreateStreamConnect(c *gin.Context) {
	userID := c.GetInt("userID")
	resultID := MustID(c, "result_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Err(err).Send()
		return
	}

	client := &Client{socket: conn, send: make(chan []byte), userID: userID}

	h.ClientManager.ResultsMap.Store(resultID, make(chan models.ResultStudent))
	ch, _ := h.ClientManager.ResultsMap.Load(resultID)

	for {
		select {
		case message := <-ch.(chan models.ResultStudent):

			err = client.socket.WriteJSON(message)

			if err != nil {
				log.Err(err).Send()
				return
			}
		}
	}
}

func (h *Handler) CreateWSStudentConnect(c *gin.Context) {
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

	passID := MustID(c, "pass_id")
	resultID := MustID(c, "result_id")

	result, err := h.TestService.GetResult(resultID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	test, err := h.TestService.Get(result.TestID, result.UserID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	questions, err := h.QuestionService.GetAllQuestionsWithAnswers(test.ID)

	if err != nil {
		SendErrorResponse(c, 401, err.Error())
		return
	}

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	timePass := 1

	h.ClientManager.PassesMap.Store(passID, make(chan models.ResultStudent))
	ch, _ := h.ClientManager.PassesMap.Load(passID)

	try := 5

	for {
		resultCh, ok := h.ClientManager.ResultsMap.Load(resultID)

		select {
		case <-time.After(time.Second * 1):
			if timePass <= result.PassageTime*60+3 {
				err := conn.WriteJSON(models.ResultStudent{
					Mark:     -1,
					TimePass: timePass,
				})

				if err != nil {
					log.Err(err).Send()

					try--

					log.Error().Int("try", try).Msg("осталось попыток")

					if try == 0 {
						resultStudent, err := h.ResultService.SaveResult(
							studentID, resultID, passID, questions, []models.QuestionWithAnswers{}, result, 0,
						)

						if err != nil {
							log.Err(err).Send()
							return
						}

						if ok {
							if resultChannel, ok := resultCh.(chan models.ResultStudent); ok {
								resultChannel <- resultStudent
							}
						}

						return
					}
				}

				if ok {
					if resultChannel, ok := resultCh.(chan models.ResultStudent); ok {
						select {
						case resultChannel <- models.ResultStudent{
							PassID:   passID,
							Mark:     -1,
							TimePass: timePass,
						}:
						default:
							log.Info().Msg("результат с каналом куда-то исчез")
						}
					}
				}

				timePass++
				log.Info().Msg("секунда отправлена")
			} else {
				log.Info().Msg("ну всё приплыли")
				resultStudent, err := h.ResultService.SaveResult(
					studentID, resultID, passID, questions, []models.QuestionWithAnswers{}, result, result.PassageTime*60,
				)

				if err != nil {
					log.Err(err).Send()
					return
				}

				if ok {
					if resultChannel, ok := resultCh.(chan models.ResultStudent); ok {
						resultChannel <- resultStudent
					}
				}

				return
			}
		case resultStudent := <-ch.(chan models.ResultStudent):
			if resultStudent.Mark != -2 {
				log.Info().Msg("тест завершён, выключаем таймер")
				return
			} else {
				log.Info().Msg("тест прерван, выключаем таймер")
				err := conn.WriteJSON(models.ResultStudent{
					Mark:     -2,
					TimePass: timePass,
				})

				if err != nil {
					log.Err(err).Send()
				}

				return
			}
		}
	}
}

func (h *Handler) ResetResult(c *gin.Context) {
	userID := c.GetInt("userID")

	resultID := MustID(c, "result_id")
	passID := MustID(c, "pass_id")

	access, err := h.TestService.GetAccess(userID, resultID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err = h.ResultService.Reset(passID, access)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	ch, ok := h.ClientManager.PassesMap.Load(passID)

	if ok {
		select {
		case ch.(chan models.ResultStudent) <- models.ResultStudent{Mark: -2}:
		default:
			log.Info().Msg("прерываем без участия")
		}
	}

	c.AbortWithStatus(204)
}
