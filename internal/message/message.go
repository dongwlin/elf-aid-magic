package message

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type Message struct {
	Type    string          `json:"type"` // request | response | event
	Action  string          `json:"action,omitempty"`
	Event   string          `json:"event,omitempty"`
	Status  string          `json:"status,omitempty"` // success | error
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data"`
	Time    time.Time       `json:"time"`
}

// Type
const (
	TypeRequest  = "request"
	TypeResponse = "response"
	TypeEvent    = "event"
)

// Status
const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func logDataMarshalError(logger *zap.Logger, err error, msgType, action, event string, data map[string]interface{}) {
	switch msgType {
	case TypeRequest, TypeResponse:
		logger.Error("Failed to marshal message data.",
			zap.Error(err),
			zap.String("type", msgType),
			zap.String("action", action),
			zap.Any("data", data),
		)
	case TypeEvent:
		logger.Error("Failed to marshal message data.",
			zap.Error(err),
			zap.String("type", msgType),
			zap.String("event", event),
			zap.Any("data", data),
		)
	default:
		logger.Error("Message type is empty or unknown.",
			zap.Error(err),
			zap.String("type", msgType),
			zap.String("action", action),
			zap.String("event", event),
			zap.Any("data", data),
		)
	}
}

func createMessage(logger *zap.Logger, msgType, action, event, status, message string, data map[string]interface{}) Message {
	if data == nil {
		data = map[string]interface{}{}
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		logDataMarshalError(logger, err, msgType, action, event, data)
		dataJSON = []byte("{}")
	}
	return Message{
		Type:    msgType,
		Action:  action,
		Event:   event,
		Status:  status,
		Message: message,
		Data:    dataJSON,
		Time:    time.Now(),
	}
}

func Createrequest(logger *zap.Logger, action string, data map[string]interface{}) Message {
	return createMessage(logger, TypeRequest, action, "", "", "", data)
}

func CreateResponse(logger *zap.Logger, action, status, message string, data map[string]interface{}) Message {
	return createMessage(logger, TypeResponse, action, "", status, message, data)
}

func CreateEvent(logger *zap.Logger, event string, data map[string]interface{}) Message {
	return createMessage(logger, TypeEvent, "", event, "", "", data)
}
