package handlers

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
}

func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
			log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				close(client.send)
				delete(manager.clients, client)
				log.Info().Msgf("Client Connected: %s", client.socket.RemoteAddr())
			}
		case message := <-manager.broadcast:
			for client := range manager.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clients, client)
				}
			}
		}
	}
}
