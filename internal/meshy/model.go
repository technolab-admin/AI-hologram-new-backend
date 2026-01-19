package meshy

// Request/Response Struct

type TextTo3DRequest struct {
	Mode         string `json:"mode"`
	Prompt       string `json:"prompt"`
	ShouldRemesh bool   `json:"shouldRemesh"`
}

type MeshyResponse struct {
	ResultID string `json:"result"`
}

type MeshyTaskStatus struct {
	Status   	string `json:"status"`
	Progress 	int    `json:"progress"`
	Mode 		string `json:"mode"`
	ModelURLS 	map[string]string `json:"model_urls"`
}
