package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

func SearchJobsJSearch(query string) ([]dtos.Job, error) {
	defaultPreference := dtos.LocationPreference{Types: []string{"remote"}}
	return SearchJobsJSearchWithLocation(query, defaultPreference)
}

func SearchJobsJSearchWithLocation(query string, locationPreference dtos.LocationPreference) ([]dtos.Job, error) {
	key := os.Getenv("RAPIDAPI_KEY")
	host := os.Getenv("RAPIDAPI_HOST")

	baseURL := fmt.Sprintf("https://%s/search", host)
	params := url.Values{}
	params.Add("query", query)
	params.Add("num_pages", "10")

	fullURL := baseURL + "?" + params.Encode()

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("X-RapidAPI-Key", key)
	req.Header.Add("X-RapidAPI-Host", host)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Data []dtos.JSearchJob `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		fmt.Println(err)
		return nil, err
	}

	jobs := []dtos.Job{}
	for _, job := range parsed.Data {
		location := "Remote"
		if !job.IsRemote {
			location = job.Location
		}
		jobs = append(jobs, dtos.Job{
			Title:    job.Title,
			Company:  job.Company,
			Location: location,
			URL:      job.URL,
			Source:   "JSearch",
		})
	}
	return jobs, nil
}
