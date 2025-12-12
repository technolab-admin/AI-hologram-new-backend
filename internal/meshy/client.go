package meshy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// This is the Meshy Client which knows how to make the POST and GET

const Meshy_API_URL string = "https://api.meshy.ai/v2/text-to-3d"

type Client struct {
	APIKey string
}

func (c *Client) getJobStatus(jobID string) (*JobStatusResponse, error) {
	url := Meshy_API_URL + jobID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var data JobStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func NewClient(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

func (c *Client) CreateGenerationJob(req *TextTo3DRequest) (*TextTo3DResponse, error) {
	body, _ := json.Marshal(req)
	fmt.Println("JSON SENT TO MESHY:", string(body))

	httpReq, _ := http.NewRequest("POST", Meshy_API_URL, bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	fmt.Println("RAW RESPONSE:", string(raw))

	var data TextTo3DResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
