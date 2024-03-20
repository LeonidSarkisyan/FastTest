package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client] = true
			log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				close(client.send)
				delete(manager.Clients, client)
				log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
			}
		case message := <-manager.Broadcast:
			for client := range manager.Clients {
				log.Info().Any("client", client).Send()
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.Clients, client)
				}
			}
		}
	}
}
