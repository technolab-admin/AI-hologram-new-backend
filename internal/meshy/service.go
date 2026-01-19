package meshy

import (
	"errors"
	"fmt"
	"log"
	"time"
	// "encoding/json"
)

// This is the service file that handles stream parsing and emits the events

type Service struct {
	client   *Client
	wsClient *WSClient
}

func NewService(client *Client, wsClient *WSClient) *Service {

	return &Service{
		client:   client,
		wsClient: wsClient,
	}
}

func (s *Service) GeneratePreview(req *TextTo3DRequest) (string, error) {

	// Generate Preview Model
	previewRes, err := s.client.CreateGenerationJob(req)
	if err != nil {
		return "", err
	}

	previewTaskID := previewRes.ResultID
	fmt.Println("Preview Task ID:" + previewTaskID)

	if err := s.waitUntilSucceeded(previewTaskID); err != nil {
		return "", err
	}

	return previewRes.ResultID, nil
}

func (s *Service) GenerateRefine(previewTaskID string) (string, error) {

	refineRes, err := s.client.CreateRefineJob(previewTaskID)
	if err != nil {
		return "", err
	}

	refineID := refineRes.ResultID

	if err := s.waitUntilSucceeded(refineID); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.glb", refineID)

	return filename, nil
}

func (s *Service) waitUntilSucceeded(taskID string) error {
	for {
		status, raw, err := s.client.getTaskStatus(taskID)
		if err != nil {
			return err
		}

		switch status.Status {

		case "SUCCEEDED":

			log.Printf("Meshy task %v succeeded", taskID) // Change to logger function

			modelName, err := download_model(raw)
			if err != nil {
				return err
			}

			err = s.wsClient.notifyFrontend(map[string]string{
				"from":   "backend-meshy",
				"target": "frontend-three",
				"event":  "new_model",
				"data":   modelName,
			})
			if err != nil {
				return err
			}

			return nil

		case "FAILED":
			return errors.New("meshy task failed") // Change to logger function
		}

		time.Sleep(2 * time.Second)

	}
}
