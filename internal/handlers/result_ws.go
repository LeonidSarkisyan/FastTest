package handlers

import (
	"App/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (h *Handler) CreateStreamConnect(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Content-Type", "text/event-stream")

	for {
		_, err := fmt.Fprintf(c.Writer, "data: %d \n\n", rand.Intn(100))
		if err != nil {
			log.Err(err).Send()
		}
		c.Writer.(http.Flusher).Flush()
		time.Sleep(1 * time.Second)
	}
}

func (h *Handler) CreateWSConnect(c *gin.Context) {
	userID := c.GetInt("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Err(err).Send()
		return
	}

	client := &Client{socket: conn, send: make(chan []byte), userID: userID}

	h.ClientManager.Register <- client
}

func (h *Handler) CreateWSStudentConnect(c *gin.Context) {
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

	client := &Client{socket: conn, send: make(chan []byte), userID: result.UserID}

	h.ClientManager.Register <- client
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

	ch, ok := h.ClientManager.ResetMap[passID]

	if !ok {
		log.Info().Msg("нет канала")
		return
	}

	*ch <- 1

	ch, ok = h.ClientManager.TimesMap[passID]

	if !ok {
		log.Info().Msg("нет канала")
		return
	}

	*ch <- 1

	h.ClientManager.Broadcast <- Message{
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
}
