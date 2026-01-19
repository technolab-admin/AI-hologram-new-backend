package websockets

import (
	"fmt"
	"slices"
	"strings"
)

var knownClients = []string{
	"backend-meshy",
	"frontend-three",
}

var requiredFields = []string{
	"from",
	"target",
	"event",
	"data",
}

func hasAllFields(msg map[string]string) error {

	for _, field := range requiredFields {
		val, exists := msg[field]
		if !exists {
			return fmt.Errorf("Missing required field: %v", field)
		}
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("Field '%v' is empty", field)
		}
	}

	return nil
}

func hasNoExtraFields(msgJson map[string]string) error {

	for field, _ := range msgJson {
		if !slices.Contains(requiredFields, field) {
			return fmt.Errorf("Unexpected field: %v", field)
		}
	}

	return nil
}

func isKnownClient(id string) bool {

	if slices.Contains(knownClients, id) {
		return true
	}

	return false
}

func vallidateClients(msg map[string]string) error {

	fromId := msg["from"]
	targetId := msg["target"]

	if !isKnownClient(fromId) {
		return fmt.Errorf("Unknown from client id: %v", fromId)
	}
	if !isKnownClient(targetId) {
		return fmt.Errorf("Unknown target client id: %v", targetId)
	}

	return nil
}

func vallidateMsg(msg map[string]string) error {

	if err := hasAllFields(msg); err != nil {
		return err
	}
	if err := hasNoExtraFields(msg); err != nil {
		return err
	}
	if err := vallidateClients(msg); err != nil {
		return err
	}

	return nil
}
