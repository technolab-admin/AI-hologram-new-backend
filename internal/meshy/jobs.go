package meshy

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
	domainErrors "AI-HOLOGRAM-NEW-BACKEND/internal/utils"
	"strconv"
	"time"
)

type JobRunner struct {
	client   *Client
	wsClient *WSClient
}

func NewJobRunner(client *Client, ws *WSClient) *JobRunner {
	return &JobRunner{client: client, wsClient: ws}
}

func (jr *JobRunner) Run(jobID string, req *TextTo3DRequest) {

	logger.Info.Println("job started: ", jobID)
	jr.send("frontend-build", "job started: ", jobID)

	// 1. Create Preview
	previewRes, err := jr.client.CreatePreviewJob(req)
	if err != nil {
		jr.fail(jobID, err)
		return
	}

	// 2. Wait for preview
	if err := jr.waitForTask(jobID, previewRes.ResultID); err != nil {
		jr.fail(jobID, err)
		return
	}

	// 3. Create Refine
	refineRes, err := jr.client.CreateRefineJob(previewRes.ResultID)
	if err != nil {
		jr.fail(jobID, err)
		return
	}

	// 4. Wait for Refine
	if err := jr.waitForTask(jobID, refineRes.ResultID); err != nil {
		jr.fail(jobID, err)
		return
	}

	jr.success(jobID)
}

func (jr *JobRunner) waitForTask(jobID, taskID string) error {

	for {
		status, raw, err := jr.client.getTaskStatus(taskID)
		if err != nil {
			return err
		}

		switch status.Status {

		case "SUCCEEDED":

			modelName, err := download_model(raw)
			if err != nil {
				return err
			}

			jr.send("frontend-three", "new_model", modelName)
			return nil

		case "FAILED":
			logger.Error.Println("task failed: ", taskID)
			return domainErrors.ErrMeshyJobFailed

		default:
			jr.send("frontend-build", "progress", strconv.Itoa(status.Progress))
		}

		time.Sleep(2 * time.Second)
	}
}

func (jr *JobRunner) send(targetID string, event string, data string) {
	_ = jr.wsClient.notifyFrontend(map[string]string{
		"from":   "backend-meshy",
		"target": targetID,
		"event":  event,
		"data":   data,
	})
}

func (jr *JobRunner) fail(jobID string, err error) {
	logger.Error.Println("job failed: ", jobID, err)

	jr.send("frontend-build", "job:error", err.Error())
}

func (jr *JobRunner) success(jobID string) {
	logger.Info.Println("job completed: ", jobID)

	jr.send("frontend-build", "job:completed", jobID)
}
