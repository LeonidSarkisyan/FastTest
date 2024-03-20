package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (h *Handler) CreateWSConnect(c *gin.Context) {
	userID := c.GetInt("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{socket: conn, send: make(chan []byte), userID: userID}

	h.ClientManager.Register <- client
}
