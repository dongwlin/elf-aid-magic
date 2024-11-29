package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dongwlin/elf-aid-magic/internal/operator"
	"github.com/gofiber/contrib/websocket"
)

type SendFunction func(msgType int, message []byte)

type WebsocketLogic struct {
	operator *operator.Operator
	sendFunc SendFunction
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWebSocketLogic(o *operator.Operator) *WebsocketLogic {
	return &WebsocketLogic{
		operator: o,
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

type Message struct {
	Type     string                 `json:"type"`
	Time     time.Time              `json:"time"`
	Payload  map[string]interface{} `json:"payload"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (l *WebsocketLogic) ProcessMessage(msg *Message) []byte {
	var resp Message

	switch msg.Type {
	case "run":
		resp = l.run(msg)
	case "stop":
		resp = l.stop(msg)
	default:
		resp = createErrorResponse(msg.Type, "Unknown message type.", nil)
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		errResponse := createErrorResponse(msg.Type, "Failed to serialize response.", nil)
		respBytes, _ = json.Marshal(errResponse)
	}

	return respBytes

}

func (l *WebsocketLogic) run(msg *Message) Message {
	if !l.operator.InitTasker() {
		l.operator.Destroy()
		return createErrorResponse(msg.Type, "Failed to init tasker.", nil)
	}

	if !l.operator.InitResource() {
		l.operator.Destroy()
		return createErrorResponse(msg.Type, "Failed to init resource.", nil)
	}

	if !l.operator.InitController() {
		l.operator.Destroy()
		return createErrorResponse(msg.Type, "Failed to init controller.", nil)
	}

	if !l.operator.Connect() {
		l.operator.Destroy()
		return createErrorResponse(msg.Type, "Failed to connect device.", nil)
	}

	ctx, cancel := context.WithCancel(context.Background())
	l.ctx = ctx
	l.cancel = cancel
	go func() {
		if l.operator.Run(l.ctx) {
			l.completed()
		}
		l.operator.Destroy()
	}()
	return createSuccessResponse(msg.Type, nil, nil)

}

func (l *WebsocketLogic) stop(msg *Message) Message {
	l.cancel()
	l.operator.Stop()
	return createSuccessResponse(msg.Type, nil, nil)
}

func (l *WebsocketLogic) completed() {
	msg := createMessage("run_completed", nil, nil)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		errResponse := createErrorResponse(msg.Type, "Failed to serialize response.", nil)
		msgBytes, _ = json.Marshal(errResponse)
	}
	l.SendBroadcastMessage(websocket.TextMessage, msgBytes)
}

func createMessage(msgType string, payload map[string]interface{}, metadata map[string]interface{}) Message {
	return Message{
		Type:     msgType,
		Time:     time.Now(),
		Payload:  payload,
		Metadata: metadata,
	}
}

func createResponse(requestType string, payload map[string]interface{}, metadata map[string]interface{}) Message {
	return createMessage(requestType+"_response", payload, metadata)
}

func createSuccessResponse(requestType string, data map[string]interface{}, metadata map[string]interface{}) Message {
	return createResponse(
		requestType,
		map[string]interface{}{
			"success": true,
			"data":    data,
		},
		metadata,
	)
}

func createErrorResponse(requestType string, err string, metadate map[string]interface{}) Message {
	return createResponse(
		requestType,
		map[string]interface{}{
			"success": false,
			"error":   err,
		},
		metadate,
	)
}
