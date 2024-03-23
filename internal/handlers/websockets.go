package handlers

import (
	"App/internal/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Message struct {
	UserID int
	PassID int
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
	passID int
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

			log.Info()

			for client := range manager.Clients {
				if client.userID != message.UserID || client.passID != message.PassID {
					continue
				}

				err := client.socket.WriteJSON(message.Result)

				if err != nil {
					log.Err(err).Send()
				}
			}
		}
	}
}
