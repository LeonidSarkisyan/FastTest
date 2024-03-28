package handlers

import (
	"App/internal/models"
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
	Clients   []*Client
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
	for _, client := range manager.Clients {
		if client.userID != message.UserID || client.passID != message.PassID {
			continue
		}
		err := client.socket.WriteJSON(message.Result)
		if err != nil {
			log.Err(err).Send()
			go manager.RemoveClient(client)
		}
	}
	manager.Mutex.Unlock()
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	manager.Clients = append(manager.Clients, client)
}

func (manager *ClientManager) RemoveClient(client *Client) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	for i, c := range manager.Clients {
		if c == client {
			manager.Clients = append(manager.Clients[:i], manager.Clients[i+1:]...)
			break
		}
	}
}
