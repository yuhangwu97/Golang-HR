package services

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket消息类型
type MessageType string

const (
	MessageTypeNotification MessageType = "notification"
	MessageTypeChat         MessageType = "chat" 
	MessageTypeSystem       MessageType = "system"
	MessageTypeHeartbeat    MessageType = "heartbeat"
)

// WebSocket消息结构
type WSMessage struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	From      string      `json:"from,omitempty"`
	To        string      `json:"to,omitempty"`
}

// 客户端连接信息
type Client struct {
	ID       string          `json:"id"`
	UserID   string          `json:"user_id"`
	Conn     *websocket.Conn `json:"-"`
	Send     chan WSMessage  `json:"-"`
	LastPing time.Time       `json:"last_ping"`
}

// WebSocket Hub 管理所有连接
type WSHub struct {
	clients    map[string]*Client
	broadcast  chan WSMessage
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// WebSocket服务
type WebSocketService struct {
	hub      *WSHub
	upgrader websocket.Upgrader
}

// 创建新的WebSocket服务
func NewWebSocketService() *WebSocketService {
	hub := &WSHub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan WSMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	return &WebSocketService{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// 在生产环境中应该检查Origin
				return true
			},
		},
	}
}

// 启动Hub
func (ws *WebSocketService) Run() {
	go ws.hub.run()
	log.Println("WebSocket service started")
}

// Hub运行主循环
func (h *WSHub) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.ID] = client
			h.mutex.Unlock()
			log.Printf("Client %s connected, total clients: %d", client.ID, len(h.clients))
			
			// 发送欢迎消息
			welcome := WSMessage{
				ID:        fmt.Sprintf("welcome_%d", time.Now().Unix()),
				Type:      MessageTypeSystem,
				Data:      map[string]string{"message": "Connected to HR Management System"},
				Timestamp: time.Now().Unix(),
			}
			select {
			case client.Send <- welcome:
			default:
				close(client.Send)
				h.mutex.Lock()
				delete(h.clients, client.ID)
				h.mutex.Unlock()
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
				log.Printf("Client %s disconnected, total clients: %d", client.ID, len(h.clients))
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mutex.RUnlock()

		case <-ticker.C:
			// 清理超时连接
			h.cleanupTimeoutClients()
		}
	}
}

// 清理超时连接
func (h *WSHub) cleanupTimeoutClients() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	timeout := time.Now().Add(-5 * time.Minute)
	for id, client := range h.clients {
		if client.LastPing.Before(timeout) {
			close(client.Send)
			delete(h.clients, id)
			log.Printf("Client %s timed out", id)
		}
	}
}

// 处理WebSocket连接
func (ws *WebSocketService) HandleConnection(c *gin.Context) {
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
	userID := c.Query("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	client := &Client{
		ID:       clientID,
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan WSMessage, 256),
		LastPing: time.Now(),
	}

	ws.hub.register <- client

	// 启动读写goroutines
	go ws.writePump(client)
	go ws.readPump(client)
}

// 读取消息
func (ws *WebSocketService) readPump(client *Client) {
	defer func() {
		ws.hub.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		client.LastPing = time.Now()
		return nil
	})

	for {
		var message WSMessage
		err := client.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message.From = client.UserID
		message.Timestamp = time.Now().Unix()
		
		// 处理不同类型的消息
		switch message.Type {
		case MessageTypeHeartbeat:
			client.LastPing = time.Now()
		default:
			ws.hub.broadcast <- message
		}
	}
}

// 发送消息
func (ws *WebSocketService) writePump(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
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

			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
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

// 广播消息给所有客户端
func (ws *WebSocketService) Broadcast(messageType MessageType, data interface{}) {
	message := WSMessage{
		ID:        fmt.Sprintf("broadcast_%d", time.Now().Unix()),
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	ws.hub.broadcast <- message
}

// 发送消息给特定用户
func (ws *WebSocketService) SendToUser(userID string, messageType MessageType, data interface{}) error {
	message := WSMessage{
		ID:        fmt.Sprintf("user_%s_%d", userID, time.Now().Unix()),
		Type:      messageType,
		Data:      data,
		To:        userID,
		Timestamp: time.Now().Unix(),
	}

	ws.hub.mutex.RLock()
	defer ws.hub.mutex.RUnlock()

	found := false
	for _, client := range ws.hub.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
				found = true
			default:
				// 客户端发送队列已满，关闭连接
				close(client.Send)
				delete(ws.hub.clients, client.ID)
			}
		}
	}

	if !found {
		return fmt.Errorf("user %s not connected", userID)
	}
	return nil
}

// 获取在线用户数量
func (ws *WebSocketService) GetOnlineUsers() int {
	ws.hub.mutex.RLock()
	defer ws.hub.mutex.RUnlock()
	return len(ws.hub.clients)
}

// 获取所有在线客户端信息
func (ws *WebSocketService) GetClients() []Client {
	ws.hub.mutex.RLock()
	defer ws.hub.mutex.RUnlock()

	clients := make([]Client, 0, len(ws.hub.clients))
	for _, client := range ws.hub.clients {
		clients = append(clients, Client{
			ID:       client.ID,
			UserID:   client.UserID,
			LastPing: client.LastPing,
		})
	}
	return clients
}

// 通知系统事件
func (ws *WebSocketService) NotifySystemEvent(event string, data map[string]interface{}) {
	notification := map[string]interface{}{
		"event":   event,
		"data":    data,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}
	ws.Broadcast(MessageTypeNotification, notification)
}