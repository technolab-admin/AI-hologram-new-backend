package meshy

import (
	"errors"
	"time"
)

// This is the service file that handles stream parsing and emits the events

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) Generate(req *TextTo3DRequest) (string, error) {

	// Create Job
	res, err := s.client.CreateGenerationJob(req)
	if err != nil {
		return "", err
	}

	// Poll Job Status (Websockets)

	for {
		status, _ := s.client.getJobStatus(res.JobID)
		if status.Status == "SUCCEEDED" {
			return status.ModelURL, nil
		}
		if status.Status == "FAILED" {
			return "", errors.New("MeshyAI generation failed")
		}
		time.Sleep(2 * time.Second)
	}
}
