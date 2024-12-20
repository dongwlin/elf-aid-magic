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
	webSocketLogic *logic.WebsocketLogic
	connections    map[*websocket.Conn]bool
	connMutex      sync.Mutex
}

func NewWebSocketHandler(logger *zap.Logger, webSocketLogic *logic.WebsocketLogic) *WebSocketHandler {
	handler := &WebSocketHandler{
		logger:         logger,
		webSocketLogic: webSocketLogic,
		connections:    make(map[*websocket.Conn]bool),
	}

	sendFunc := func(msgType int, message []byte) {
		handler.BroadcastMessage(msgType, message)
	}

	webSocketLogic.SetSendFunction(sendFunc)
	return handler
}

func (h *WebSocketHandler) Register(r fiber.Router) {
	r.Get("/", websocket.New(h.WebSocket))
}

func (h *WebSocketHandler) WebSocket(c *websocket.Conn) {
	h.connMutex.Lock()
	h.connections[c] = true
	h.connMutex.Unlock()

	defer func() {
		h.connMutex.Lock()
		delete(h.connections, c)
		h.connMutex.Unlock()
		c.Close()
	}()

	var (
		mt     int
		msgBuf []byte
		err    error
	)
	for {
		if mt, msgBuf, err = c.ReadMessage(); err != nil {
			h.logger.Error(
				"failed to read msg",
				zap.Error(err),
			)
			break
		}

		var msg message.Message
		err = json.Unmarshal(msgBuf, &msg)
		if err != nil {
			h.logger.Error(
				"failed to unmarshal msg",
				zap.Error(err),
			)
			break
		}
		h.logger.Info(
			"recv msg",
			zap.String("type", msg.Type),
			zap.Time("time", msg.Time),
			zap.Any("payload", msg.Data),
		)

		resp := h.webSocketLogic.ProcessMessage(&msg)

		if err = c.WriteMessage(mt, resp); err != nil {
			h.logger.Error(
				"failed to write msg",
				zap.Error(err),
			)
			break
		}
	}

}

func (h *WebSocketHandler) BroadcastMessage(msgType int, message []byte) {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	for conn := range h.connections {
		if err := conn.WriteMessage(msgType, message); err != nil {
			h.logger.Error("failed to send message", zap.Error(err))
			conn.Close()
			delete(h.connections, conn)
		}
	}
}
