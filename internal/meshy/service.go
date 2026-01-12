package meshy

import (
	"errors"
	"fmt"
	"time"
)

// This is the service file that handles stream parsing and emits the events

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
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

	return refineRes.ResultID, nil
}

func (s *Service) waitUntilSucceeded(taskID string) error {
	for {
		status, err := s.client.getTaskStatus(taskID)
		if err != nil {
			return err
		}

		switch status.Status {
		case "SUCCEEDED":
			return nil
		case "FAILED":
			return errors.New("meshy task failed")
		}

		time.Sleep(2 * time.Second)

	}
}
