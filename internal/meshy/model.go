package meshy

// Request/Response Struct

type TextTo3DRequest struct {
	Mode         string `json:"mode"`
	Prompt       string `json:"prompt"`
	ShouldRemesh bool   `json:"shouldRemesh"`
}
type TextTo3DResponse struct {
	JobID string `json:"job_id"`
}

type JobStatusResponse struct {
	Status   string `json:"status"`
	Progress int    `json:progress`
	ModelURL string `json:"model_url"`
	Error    string `json:error,omitempty`
}
