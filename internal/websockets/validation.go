package websockets

import (
	"fmt"
	"slices"
	"strings"
)

var knownClients = []string{
	"backend-meshy",
	"frontend-three",
	"frontend-build",

	"TEST-A",
	"TEST-B",
}

var requiredFields = []string{
	"from",
	"target",
	"event",
	"data",
}


func HasAllFields(msg map[string]string) error {

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


func HasNoExtraFields(msgJson map[string]string) error {

	for field, _ := range msgJson {
		if !slices.Contains(requiredFields, field) {
			return fmt.Errorf("Unexpected field: %v", field)
		}
	}

	return nil
}

func IsKnownClient(id string) bool {

	if slices.Contains(knownClients, id) {
		return true
	}

	return false
}

func VallidateClients(msg map[string]string) error {

	fromId := msg["from"]
	targetId := msg["target"]

	if !IsKnownClient(fromId) {
		return fmt.Errorf("Unknown from client id: %v", fromId)
	}
	if !IsKnownClient(targetId) {
		return fmt.Errorf("Unknown target client id: %v", targetId)
	}

	return nil
}

func VallidateMsg(msg map[string]string) error {

	if err := HasAllFields(msg); err != nil {
		return err
	}
	if err := HasNoExtraFields(msg); err != nil {
		return err
	}
	if err := VallidateClients(msg); err != nil {
		return err
	}

	return nil
}

