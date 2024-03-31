package handlers

import (
	"App/internal/models"
	"github.com/gorilla/websocket"
	"sync"
)

type Message struct {
	UserID int
	PassID int
	Result models.ResultStudent
}

type ClientManager struct {
	Clients    sync.Map
	TimesMap   sync.Map
	ResetMap   sync.Map
	ResultsMap sync.Map
	PassesMap  sync.Map
	Mutex      sync.Mutex
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	userID int
	passID int
}
