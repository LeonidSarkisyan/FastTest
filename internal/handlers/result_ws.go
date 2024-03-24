package handlers

import (
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Err(err).Send()
		return
	}

	client := &Client{socket: conn, send: make(chan []byte), userID: userID}

	h.ClientManager.AddClient(client)
}

func (h *Handler) CreateWSStudentConnect(c *gin.Context) {
	passID := MustID(c, "pass_id")
	resultID := MustID(c, "result_id")

	result, err := h.TestService.GetResult(resultID)

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

	client := &Client{socket: conn, send: make(chan []byte), userID: result.UserID, passID: passID}

	h.ClientManager.AddClient(client)
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

	go func() {
		message := Message{
			UserID: userID,
			Result: models.ResultStudent{
				Mark:         -2,
				Score:        0,
				MaxScore:     0,
				DateTimePass: time.Time{},
				PassID:       passID,
				AccessID:     0,
				StudentID:    0,
				TimePass:     0,
			},
		}

		h.SendToBroadcast(message)
		message.PassID = passID
		h.SendToBroadcast(message)

		h.ClientManager.TimesMap[passID] <- 1
	}()
}
