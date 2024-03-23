package handlers

import (
	"App/internal/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Message struct {
	UserID int
	Result models.ResultStudent
}

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	TimesMap   map[int]*chan int
	ResetMap   map[int]*chan int
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	userID int
}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client] = true
			log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
			log.Info().Any("manager clients", manager.Clients).Send()
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				close(client.send)
				delete(manager.Clients, client)
				log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
			}
		case message := <-manager.Broadcast:
			for client := range manager.Clients {
				log.Info().Int("client user_id", client.userID).Send()
				log.Info().Int("user_id", message.UserID).Send()

				if client.userID != message.UserID {
					continue
				}

				resultMessage, err := json.Marshal(message.Result)

				if err != nil {
					log.Err(err).Send()
					return
				}

				err = client.socket.WriteMessage(websocket.TextMessage, resultMessage)

				if err != nil {
					log.Err(err).Send()
				}
			}
		}
	}
}
