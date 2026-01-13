package meshy


import (
	"encoding/json"
	"path/filepath"
	"io"
	"net/http"
	"os"
)


func download_model(raw []byte) (string, error) {

	var status MeshyTaskStatus
	if err := json.Unmarshal(raw, &status); err != nil {
		return "", err
	}

	url := status.ModelURLS["glb"]
	modelName := status.Mode

	outPath := filepath.Join("assets", "downloads", modelName+".glb")

	resp, err := http.Get(url)
	if err != nil {return "", err}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {return "",err}

	outFile, err := os.Create(outPath)
	if err != nil {return "", err}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {return "", err}

	return modelName, nil
}