package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSClient struct {
	UserID int
	Conn   *websocket.Conn
	Send   chan []byte
}

type WSHub struct {
	clients    map[int]*WSClient
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	mutex      sync.RWMutex
}

var Hub = &WSHub{
	clients:    make(map[int]*WSClient),
	broadcast:  make(chan []byte),
	register:   make(chan *WSClient),
	unregister: make(chan *WSClient),
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.UserID] = client
			h.mutex.Unlock()
			log.Printf("WebSocket: User %d connected, total clients: %d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mutex.Unlock()
			log.Printf("WebSocket: User %d disconnected, total clients: %d", client.UserID, len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for userID, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, userID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *WSHub) SendToUser(userID int, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, userID)
		}
	}
}

type WSMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type NotificationPayload struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type MessagePayload struct {
	FromUserID   int    `json:"from_user_id"`
	FromUsername string `json:"from_username"`
	Content      string `json:"content"`
	SentAt       string `json:"sent_at"`
}

func init() {
	go Hub.Run()
}

type WSHandler struct{}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	Hub.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHandler) readPump(client *WSClient) {
	defer func() {
		Hub.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err == nil {
			log.Printf("WebSocket message from user %d: type=%s", client.UserID, msg.Type)
			handleClientMessage(client, msg)
		}
	}
}

func (h *WSHandler) writePump(client *WSClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleClientMessage(client *WSClient, msg WSMessage) {
	switch msg.Type {
	case "ping":
		sendMessage(client, WSMessage{Type: "pong", Content: "ok"})
	case "subscribe":
		sendMessage(client, WSMessage{Type: "subscribed", Content: map[string]interface{}{
			"user_id":  client.UserID,
			"status":   "connected",
			"timestamp": time.Now().Unix(),
		}})
	}
}

func sendMessage(client *WSClient, msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case client.Send <- data:
	default:
	}
}

func SendNotificationToUser(userID int, notification NotificationPayload) {
	msg := WSMessage{
		Type:    "notification",
		Content: notification,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	Hub.SendToUser(userID, data)
}

func SendMessageToUser(userID int, message MessagePayload) {
	msg := WSMessage{
		Type:    "new_message",
		Content: message,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	Hub.SendToUser(userID, data)
}
