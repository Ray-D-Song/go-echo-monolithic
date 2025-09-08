package service

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// WebSocketService handles WebSocket connections
type WebSocketService struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]*Client
	mutex    sync.RWMutex
	logger   *zap.Logger
}

// Client represents a WebSocket client
type Client struct {
	conn   *websocket.Conn
	userID uint
	send   chan []byte
}

// Message represents a WebSocket message
type Message struct {
	Type    string `json:"type"`
	Data    any    `json:"data"`
	UserID  uint   `json:"user_id,omitempty"`
	Channel string `json:"channel,omitempty"`
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(logger *zap.Logger) *WebSocketService {
	return &WebSocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin
				// In production, you should validate the origin
				return true
			},
		},
		clients: make(map[*websocket.Conn]*Client),
		logger:  logger,
	}
}

// UpgradeConnection upgrades HTTP connection to WebSocket
func (s *WebSocketService) UpgradeConnection(c echo.Context, userID uint) error {
	conn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		s.logger.Error("Failed to upgrade connection", zap.Error(err))
		return err
	}

	client := &Client{
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
	}

	s.mutex.Lock()
	s.clients[conn] = client
	s.mutex.Unlock()

	s.logger.Info("Client connected", zap.Uint("user_id", userID))

	// Start goroutines for reading and writing
	go s.readPump(client)
	go s.writePump(client)

	return nil
}

// readPump handles reading messages from the WebSocket connection
func (s *WebSocketService) readPump(client *Client) {
	defer func() {
		s.removeClient(client.conn)
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Time{})
		return nil
	})

	for {
		var message Message
		err := client.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error("Unexpected WebSocket close", zap.Error(err))
			}
			break
		}

		// Handle different message types
		s.handleMessage(client, &message)
	}
}

// writePump handles writing messages to the WebSocket connection
func (s *WebSocketService) writePump(client *Client) {
	defer client.conn.Close()

	for message := range client.send {
		client.conn.SetWriteDeadline(time.Time{})
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			s.logger.Error("Failed to write message", zap.Error(err))
			return
		}
	}

	client.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// handleMessage processes incoming WebSocket messages
func (s *WebSocketService) handleMessage(client *Client, message *Message) {
	s.logger.Info("Received message",
		zap.String("type", message.Type),
		zap.Uint("user_id", client.userID),
	)

	switch message.Type {
	case "ping":
		s.sendToClient(client, &Message{
			Type: "pong",
			Data: "pong",
		})
	case "echo":
		s.sendToClient(client, &Message{
			Type: "echo",
			Data: message.Data,
		})
	case "broadcast":
		s.broadcastMessage(message, client.userID)
	default:
		s.logger.Warn("Unknown message type", zap.String("type", message.Type))
	}
}

// sendToClient sends a message to a specific client
func (s *WebSocketService) sendToClient(client *Client, message *Message) {
	select {
	case client.send <- s.encodeMessage(message):
	default:
		close(client.send)
		s.removeClient(client.conn)
	}
}

// broadcastMessage sends a message to all connected clients
func (s *WebSocketService) broadcastMessage(message *Message, senderID uint) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		if client.userID != senderID { // Don't send to sender
			select {
			case client.send <- s.encodeMessage(message):
			default:
				close(client.send)
				delete(s.clients, client.conn)
			}
		}
	}
}

// SendToUser sends a message to a specific user
func (s *WebSocketService) SendToUser(userID uint, message *Message) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		if client.userID == userID {
			select {
			case client.send <- s.encodeMessage(message):
			default:
				close(client.send)
				delete(s.clients, client.conn)
			}
		}
	}
}

// removeClient removes a client from the clients map
func (s *WebSocketService) removeClient(conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[conn]; exists {
		close(client.send)
		delete(s.clients, conn)
		s.logger.Info("Client disconnected", zap.Uint("user_id", client.userID))
	}
}

// encodeMessage encodes a message to JSON bytes
func (s *WebSocketService) encodeMessage(message *Message) []byte {
	// In a real implementation, you would use json.Marshal
	// For simplicity, returning a placeholder
	return []byte(`{"type":"` + message.Type + `","data":"` + message.Type + `"}`)
}

// GetConnectedUsers returns the number of connected users
func (s *WebSocketService) GetConnectedUsers() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}
