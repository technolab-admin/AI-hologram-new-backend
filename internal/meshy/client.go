package meshy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// This is the Meshy Client which knows how to make the POST and GET

type Client struct {
	APIKey  string
	BaseURL string
	http    *http.Client
}

func NewClient(apiKey string, baseURL string) *Client {
	return &Client{APIKey: apiKey, BaseURL: baseURL, http: &http.Client{}}
}

func (c *Client) CreateGenerationJob(req *TextTo3DRequest) (*MeshyResponse, error) {
	url := c.BaseURL + "/text-to-3d"
	return c.CreateJob(url, req)
}

func (c *Client) CreateRefineJob(previewID string) (*MeshyResponse, error) {
	url := c.BaseURL + "/text-to-3d"
	req := map[string]string{
		"mode":            "refine",
		"preview_task_id": previewID,
	}
	return c.CreateJob(url, req)
}

func (c *Client) CreateJob(url string, payload any) (*MeshyResponse, error) {
	body, _ := json.Marshal(payload)
	fmt.Println("JSON SENT TO MESHY:", string(body)) // debug

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	fmt.Println("RAW RESPONSE:", string(raw)) // debug

	var data MeshyResponse
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) getTaskStatus(taskID string) (*MeshyTaskStatus, []byte, error) {
	url := c.BaseURL + "/text-to-3d/" + taskID

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	fmt.Println("TASK STATUS RAW: ", string(raw)) // debug

	var status MeshyTaskStatus
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, raw, err
	}

	return &status, raw, nil
}
