package handler

import (
	"encoding/json"
	"sync"

	"github.com/dongwlin/elf-aid-magic/internal/logic"
	"github.com/dongwlin/elf-aid-magic/internal/message"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WebSocketHandler struct {
	logger         *zap.Logger
	webSocketLogic *logic.WebSocketLogic
	connections    map[*websocket.Conn]bool
	connMutex      sync.Mutex
}

func NewWebSocketHandler(logger *zap.Logger, webSocketLogic *logic.WebSocketLogic) *WebSocketHandler {
	handler := &WebSocketHandler{
		logger:         logger,
		webSocketLogic: webSocketLogic,
		connections:    make(map[*websocket.Conn]bool),
	}

	webSocketLogic.SetSendMessageFunc(handler.SendMessage)
	webSocketLogic.SetBroadcastMessageFunc(handler.BroadcastMessage)

	return handler
}

func (h *WebSocketHandler) Register(r fiber.Router) {
	r.Get("/", websocket.New(h.WebSocket))
}

func (h *WebSocketHandler) WebSocket(c *websocket.Conn) {
	h.addConnection(c)
	defer h.removeConnection(c)

	for {
		if err := h.handleConnection(c); err != nil {
			h.logger.Error("connection error",
				zap.Error(err),
			)
			break
		}
	}

}

func (h *WebSocketHandler) addConnection(c *websocket.Conn) {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()
	h.connections[c] = true
}

func (h *WebSocketHandler) removeConnection(c *websocket.Conn) {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()
	delete(h.connections, c)
	c.Close()
}

func (h *WebSocketHandler) handleConnection(c *websocket.Conn) error {
	msgType, msgBuf, err := c.ReadMessage()
	if err != nil {
		h.logger.Error("failed to read message",
			zap.Error(err),
		)
		return err
	}

	var msg message.Message
	err = json.Unmarshal(msgBuf, &msg)
	if err != nil {
		h.logger.Error("failed to unmarshal message",
			zap.Error(err),
		)
		return err
	}

	h.webSocketLogic.ProcessMessage(c, msgType, &msg)

	return nil
}

func (h *WebSocketHandler) SendMessage(conn *websocket.Conn, msgType int, message []byte) error {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	if _, exists := h.connections[conn]; exists {
		if err := conn.WriteMessage(msgType, message); err != nil {
			h.logger.Error("Failed to send message.",
				zap.Error(err),
			)
			conn.Close()
			delete(h.connections, conn)
			return err
		}
	}
	return nil
}

func (h *WebSocketHandler) BroadcastMessage(msgType int, message []byte) {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	for conn := range h.connections {
		if err := conn.WriteMessage(msgType, message); err != nil {
			h.logger.Error("Failed to send message by broadcast.",
				zap.Error(err),
			)
			conn.Close()
			delete(h.connections, conn)
		}
	}
}
