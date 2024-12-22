package logic

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dongwlin/elf-aid-magic/internal/message"
	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

type SendMessageFunc func(conn *websocket.Conn, msgType int, msg []byte) error

type BroadcastMessageFunc func(msgType int, msg []byte)

type WebsocketLogic struct {
	logger               *zap.Logger
	operatorManager      *operator.Manager
	sendMessageFunc      SendMessageFunc
	broadcastMessageFunc BroadcastMessageFunc
	ctx                  context.Context
	cancel               context.CancelFunc
}

func NewWebSocketLogic(logger *zap.Logger, om *operator.Manager) *WebsocketLogic {
	return &WebsocketLogic{
		logger:          logger,
		operatorManager: om,
	}
}

func (l *WebsocketLogic) SetSendMessageFunc(sendMessageFunc SendMessageFunc) {
	l.sendMessageFunc = sendMessageFunc
}

func (l *WebsocketLogic) sendMessage(conn *websocket.Conn, msgType int, message []byte) error {
	if l.sendMessageFunc != nil {
		return l.sendMessageFunc(conn, msgType, message)
	}
	return errors.New("WebSocketLogic sendMessageFunc is nil")
}

func (l *WebsocketLogic) SetBroadcastMessageFunc(broadcastMessageFunc BroadcastMessageFunc) {
	l.broadcastMessageFunc = broadcastMessageFunc
}

func (l *WebsocketLogic) broadcastMessage(msgType int, message []byte) {
	if l.broadcastMessageFunc != nil {
		l.broadcastMessageFunc(msgType, message)
	}
}

func (l *WebsocketLogic) ProcessMessage(conn *websocket.Conn, msgType int, msg *message.Message) {
	switch msg.Type {
	case message.TypeRequest:
		l.handleRequest(conn, msgType, msg)
	case message.TypeResponse:
		l.handleResponse(conn, msgType, msg)
	case message.TypeEvent:
		l.handleEvent(conn, msgType, msg)
	default:
		l.handleUnknowMessageType(conn, msgType, msg)
	}
}

func (l *WebsocketLogic) handleRequest(conn *websocket.Conn, msgType int, msg *message.Message) {
	var resp message.Message
	switch msg.Action {
	case "start":
		resp = l.start(msg)
	case "stop":
		resp = l.stop(msg)
	default:
		resp = message.CreateResponse(l.logger, msg.Action, message.StatusError, "Unknown request action.", nil)
	}

	respBytes := serializeMessage(l.logger, resp)
	err := l.sendMessage(conn, msgType, respBytes)
	if err != nil {
		l.logger.Error("failed to send message",
			zap.Error(err),
		)
	}
}

func (l *WebsocketLogic) handleResponse(conn *websocket.Conn, msgType int, msg *message.Message) {
	l.logger.Error("unknown response action",
		zap.String("action", msg.Action),
	)
	event := message.CreateEvent(l.logger, "UnknownResponseAction", nil)
	eventBytes := serializeMessage(l.logger, event)
	err := l.sendMessage(conn, msgType, eventBytes)
	if err != nil {
		l.logger.Error("failed to send message",
			zap.Error(err),
		)
	}
}

func (l *WebsocketLogic) handleEvent(conn *websocket.Conn, msgType int, msg *message.Message) {
	l.logger.Error("unknown event",
		zap.String("event", msg.Event),
	)
	event := message.CreateEvent(l.logger, "UnknownEvent", nil)
	eventBytes := serializeMessage(l.logger, event)
	err := l.sendMessage(conn, msgType, eventBytes)
	if err != nil {
		l.logger.Error("failed to send message",
			zap.Error(err),
		)
	}
}

func (l *WebsocketLogic) handleUnknowMessageType(conn *websocket.Conn, msgType int, msg *message.Message) {
	l.logger.Error("unknown message type",
		zap.String("type", msg.Type),
	)
	event := message.CreateEvent(l.logger, "UnknownMessageType", nil)
	eventBytes := serializeMessage(l.logger, event)
	err := l.sendMessage(conn, msgType, eventBytes)
	if err != nil {
		l.logger.Error("failed to send message",
			zap.Error(err),
		)
	}
}

func serializeMessage(logger *zap.Logger, msg message.Message) []byte {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		errEvent := message.CreateEvent(logger, "SerializeMessageError", nil)
		msgBytes, _ = json.Marshal(errEvent)
	}
	return msgBytes
}

type MessageStartRequestData struct {
	TaskerID string `json:"tasker_id"`
}

func (l *WebsocketLogic) start(msg *message.Message) message.Message {
	var data MessageStartRequestData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to unserialize request data.", nil)
	}

	if data.TaskerID == "" {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Tasker ID is empty.", nil)
	}

	operator, exists := l.operatorManager.GetOperatorByID(data.TaskerID)
	if !exists {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Operator don't exists.", nil)
	}
	if !operator.InitTasker() {
		operator.Destroy()
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to init tasker.", nil)
	}

	if !operator.InitResource() {
		operator.Destroy()
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to init resource.", nil)
	}

	if !operator.InitController() {
		operator.Destroy()
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to init controller.", nil)
	}

	if !operator.Connect() {
		operator.Destroy()
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to connect device.", nil)
	}

	ctx, cancel := context.WithCancel(context.Background())
	l.ctx = ctx
	l.cancel = cancel
	go func() {
		if operator.Run(l.ctx) {
			l.completed(operator.ID)
		}
		operator.Destroy()
	}()
	return message.CreateResponse(l.logger, msg.Action, message.StatusSuccess, "Success", nil)

}

type MessageStopRequestData struct {
	TaskerID string `json:"tasker_id"`
}

func (l *WebsocketLogic) stop(msg *message.Message) message.Message {
	l.cancel()

	var data MessageStopRequestData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to unserialize request data.", nil)
	}

	if data.TaskerID == "" {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Tasker ID is empty.", nil)
	}

	operator, exists := l.operatorManager.GetOperatorByID(data.TaskerID)
	if !exists {
		return message.CreateResponse(l.logger, msg.Action, message.StatusError, "Operator don't exists.", nil)
	}
	operator.Stop().Wait()
	return message.CreateResponse(l.logger, msg.Action, message.StatusSuccess, "Success", nil)
}

type EventMessageCompletedData struct {
	TaskerID string `json:"tasker_id"`
}

func (l *WebsocketLogic) completed(taskerID string) {
	data := EventMessageCompletedData{
		TaskerID: taskerID,
	}
	msg := message.CreateEvent(l.logger, "completed", data)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		errResponse := message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to serialize response.", nil)
		msgBytes, _ = json.Marshal(errResponse)
	}
	l.broadcastMessage(websocket.TextMessage, msgBytes)
}
