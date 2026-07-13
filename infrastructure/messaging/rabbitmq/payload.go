package rabbitmq

import (
	"encoding/json"
	"fmt"
)

type envelope struct {
	SecretKey string          `json:"secret_key"`
	Message   json.RawMessage `json:"message"`
}

func UnmarshalPayload(body []byte, dest any) (string, error) {
	var payload envelope
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	if err := json.Unmarshal(payload.Message, dest); err != nil {
		return "", fmt.Errorf("invalid message payload: %w", err)
	}

	return payload.SecretKey, nil
}
