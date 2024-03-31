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
	Clients  sync.Map
	TimesMap sync.Map
	ResetMap sync.Map
	Mutex    sync.Mutex
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	userID int
	passID int
}

func (manager *ClientManager) SendToBroadcast(message Message) {
	manager.Clients.Range(func(key, value interface{}) bool {
		client := key.(*Client)

		if client.userID != message.UserID || client.passID != message.PassID {
			return true
		}

		manager.Mutex.Lock()
		err := client.socket.WriteJSON(message.Result)
		manager.Mutex.Unlock()

		if err != nil {
			log.Err(err).Send()
			manager.RemoveClient(client)
		}

		return true
	})
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.Clients.Store(client, true)
}

func (manager *ClientManager) RemoveClient(client *Client) {
	manager.Clients.Delete(client)
}
