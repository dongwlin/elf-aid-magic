package message

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// Message represents the structure of a WebSocket message used in the application.
// It supports both request-response and event-driven communication patterns.
//
// Fields for request-response messages:
// - Type:    The type of message, either "request" or "response".
// - Action:  Specifies the action to be performed (e.g., "startTask").
// - Status:  Indicates the status of the response ("success" or "error").
// - Message: Optional descriptive information, such as error details or success messages.
// - Data:    Contains the payload of the request or response.
// - Time:    Timestamp of when the message was created.
//
// Fields for event messages:
// - Type:    The type of message, always "event".
// - Event:   Specifies the name of the event (e.g., "taskProgress").
// - Data:    Contains the event payload.
// - Time:    Timestamp of when the event was created.
//
// Notes:
// - The "Action", "Status" and "Message" fields are omitted for event messages.
// - The "Event" field is omitted for request-response messages.
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
	fields := []zap.Field{
		zap.Error(err),
		zap.String("type", msgType),
		zap.Any("data", data),
	}

	if msgType == TypeRequest || msgType == TypeResponse {
		fields = append(fields, zap.String("action", action))
	} else if msgType == TypeEvent {
		fields = append(fields, zap.String("event", event))
	} else {
		fields = append(fields, zap.String("action", action), zap.String("event", event))
	}

	logger.Error("Failed to marshal message data.", fields...)
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

func CreateRequest(logger *zap.Logger, action string, data map[string]interface{}) Message {
	return createMessage(logger, TypeRequest, action, "", "", "", data)
}

func CreateResponse(logger *zap.Logger, action, status, message string, data map[string]interface{}) Message {
	return createMessage(logger, TypeResponse, action, "", status, message, data)
}

func CreateEvent(logger *zap.Logger, event string, data map[string]interface{}) Message {
	return createMessage(logger, TypeEvent, "", event, "", "", data)
}
