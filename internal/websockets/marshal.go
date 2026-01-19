package websockets

import (
	"encoding/json"
	"fmt"
)

func MarshalMessage(msg map[string]string) ([]byte, error) {

	rawJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Error converting message to bytes: %v", err) // Change to logger function
	}
	return rawJSON, nil
}

func UnmarshalMessage(raw []byte) (map[string]string, error) {

	var msg map[string]string
	err := json.Unmarshal(raw, &msg)
	if err != nil {
		return nil, fmt.Errorf("Error converting message to JSON: %v", err) // Change to logger function
	}
	return msg, err
}
