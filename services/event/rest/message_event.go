package event_rest

import (
	"encoding/json"
)

type MessageEventObject struct {
	Message string `json:"message"`
}

func CreateMessageEventObject(message string) ([]byte, error) {
	messageObject := MessageEventObject{
		Message: message,
	}
	result, err := json.Marshal(messageObject)
	if err != nil {
		return nil, err
	}
	return result, nil
}
