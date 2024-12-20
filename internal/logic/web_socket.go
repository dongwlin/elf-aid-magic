package logic

import (
	"context"
	"encoding/json"

	"github.com/dongwlin/elf-aid-magic/internal/message"
	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

type SendFunction func(msgType int, message []byte)

type WebsocketLogic struct {
	logger          *zap.Logger
	operatorManager *operator.Manager
	sendFunc        SendFunction
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewWebSocketLogic(logger *zap.Logger, om *operator.Manager) *WebsocketLogic {
	return &WebsocketLogic{
		logger:          logger,
		operatorManager: om,
	}
}

func (l *WebsocketLogic) SetSendFunction(sendFunc SendFunction) {
	l.sendFunc = sendFunc
}

func (l *WebsocketLogic) SendBroadcastMessage(msgType int, message []byte) {
	if l.sendFunc != nil {
		l.sendFunc(msgType, message)
	}
}

func (l *WebsocketLogic) ProcessMessage(msg *message.Message) []byte {
	var resp message.Message

	switch msg.Type {
	case "run":
		resp = l.run(msg)
	case "stop":
		resp = l.stop(msg)
	default:
		resp = message.CreateResponse(l.logger, msg.Action, message.StatusError, "Unknown message action.", nil)
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		errResponse := message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to serialize response.", nil)
		respBytes, _ = json.Marshal(errResponse)
	}

	return respBytes

}

type MessageRunRequestData struct {
	TaskerID string `json:"tasker_id"`
}

func (l *WebsocketLogic) run(msg *message.Message) message.Message {
	var data MessageRunRequestData
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
			l.completed()
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

func (l *WebsocketLogic) completed() {
	msg := message.CreateEvent(l.logger, "run_completed", nil)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		errResponse := message.CreateResponse(l.logger, msg.Action, message.StatusError, "Failed to serialize response.", nil)
		msgBytes, _ = json.Marshal(errResponse)
	}
	l.SendBroadcastMessage(websocket.TextMessage, msgBytes)
}
