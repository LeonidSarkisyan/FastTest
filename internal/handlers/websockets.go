package handlers

import (
	"App/internal/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"sync"
)

type Message struct {
	UserID int
	PassID int
	Result models.ResultStudent
}

type ClientManager struct {
	Clients   map[*Client]struct{}
	Broadcast chan Message
	TimesMap  map[int]chan int
	ResetMap  map[int]chan int
	Mutex     sync.Mutex
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	userID int
	passID int
}

func (manager *ClientManager) SendToBroadcast(message Message) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	for client := range manager.Clients {
		if client.userID != message.UserID || client.passID != message.PassID {
			continue
		}

		jsonData, err := json.Marshal(message.Result)
		if err != nil {
			log.Error().Err(err).Msg("Error marshaling message")
			continue
		}

		select {
		case client.send <- jsonData:
		default:
			log.Error().Msg("Failed to send message to client")
			manager.RemoveClient(client)
		}
	}
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	manager.Clients[client] = struct{}{}
}

func (manager *ClientManager) RemoveClient(client *Client) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	delete(manager.Clients, client)
}

func (client *Client) Read() {
	defer func() {
		client.socket.Close()
	}()

	for {
		_, _, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (client *Client) Write() {
	defer func() {
		client.socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}
			err := client.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}
