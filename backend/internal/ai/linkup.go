package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func SearchJobsLinkUp(query string) ([]dtos.Job, error) {
	apiKey := os.Getenv("LINKUP_API_KEY")
	url := os.Getenv("LINKUP_API_URL")

	payload := map[string]interface{}{
		"q":             query,
		"depth":         "standard",
		"outputType":    "searchResults",
		"includeImages": false,
	}

	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Results []dtos.LinkupJob `json:"results"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	jobs := []dtos.Job{}
	for _, r := range parsed.Results {
		jobs = append(jobs, dtos.Job{
			Title:       r.Name,
			Description: r.Content,
			URL:         r.URL,
			Source:      "LinkUp",
		})
	}

	return jobs, nil
}
