package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func SearchJobsLinkUpStructured(query string) (*dtos.JobAnnouncements, error) {
	defaultPreference := dtos.LocationPreference{Types: []string{"remote"}}
	return SearchJobsLinkUpStructuredWithLocation(query, defaultPreference)
}

func SearchJobsLinkUpStructuredWithLocation(query string, locationPreference dtos.LocationPreference) (*dtos.JobAnnouncements, error) {
	apiKey := os.Getenv("LINKUP_API_KEY")
	url := os.Getenv("LINKUP_API_URL")

	if apiKey == "" || url == "" {
		return nil, fmt.Errorf("LINKUP_API_KEY and LINKUP_API_URL environment variables are required")
	}

	structuredSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"jobs": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"job_title": map[string]interface{}{
							"type":        "string",
							"description": "Job Title mentioned in the job announcement",
						},
						"experience_level": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"internship", "entry level", "junior", "mid-level", "senior"},
							"description": "Required experience level",
						},
						"required_skills": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of required skills for the job",
						},
						"remote": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the job is remote or not",
						},
						"location": map[string]interface{}{
							"type":        "string",
							"description": "Location, if there is any location restriction in the job",
						},
						"salary": map[string]interface{}{
							"type":        "integer",
							"description": "Yearly salary, when available",
						},
						"job_post_url": map[string]interface{}{
							"type":        "string",
							"description": "URL to the job announcement",
						},
						"company": map[string]interface{}{
							"type":        "string",
							"description": "Company hiring for the job",
						},
					},
					"required": []string{"job_title", "experience_level", "required_skills", "remote", "job_post_url", "company"},
				},
			},
		},
		"required": []string{"jobs"},
	}

	payload := map[string]interface{}{
		"q":                      query,
		"depth":                  "standard",
		"outputType":             "structured",
		"includeImages":          false,
		"structuredOutputSchema": structuredSchema,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var jobAnnouncements dtos.JobAnnouncements
	if err := json.Unmarshal(respBody, &jobAnnouncements); err != nil {
		var oldFormat struct {
			Results []dtos.LinkupJob `json:"results"`
		}
		if fallbackErr := json.Unmarshal(respBody, &oldFormat); fallbackErr != nil {
			return nil, fmt.Errorf("failed to parse structured response: %w (original error: %v)", fallbackErr, err)
		}

		jobAnnouncements = convertLinkupJobsToStructured(oldFormat.Results)
	}

	return &jobAnnouncements, nil
}

func convertLinkupJobsToStructured(linkupJobs []dtos.LinkupJob) dtos.JobAnnouncements {
	var jobs []dtos.JobDescription

	for _, linkupJob := range linkupJobs {
		job := dtos.JobDescription{
			JobTitle:        linkupJob.Name,
			ExperienceLevel: "mid-level",
			RequiredSkills:  []string{},
			Remote:          false,
			Location:        nil,
			Salary:          nil,
			JobPostURL:      linkupJob.URL,
			Company:         "Unknown",
		}
		jobs = append(jobs, job)
	}

	return dtos.JobAnnouncements{Jobs: jobs}
}
