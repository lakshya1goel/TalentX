package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
	"google.golang.org/genai"
)

type AIClient struct {
	Client *genai.Client
}

func NewAIClient(ctx context.Context, apiKey string) *AIClient {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		fmt.Println(err)
	}

	return &AIClient{Client: client}
}

func (a *AIClient) GetJobsFromResume(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.Job, error) {
	parsedLocationPreference := a.parseLocationPreference(locationPreference)
	prompt := a.Prompt(parsedLocationPreference)

	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: "application/pdf",
				Data:     pdfBytes,
			},
		},
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	tools := a.Tools()

	result, err := a.Client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Tools: tools,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	var functionCalls []*genai.FunctionCall
	for _, candidate := range result.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.FunctionCall != nil {
				functionCalls = append(functionCalls, part.FunctionCall)
				fmt.Printf("AI wants to call: %s\n", part.FunctionCall.Name)
			}
		}
	}

	if len(functionCalls) == 0 {
		fmt.Println("No function calls found in AI response")
		return []dtos.Job{}, nil
	}

	return a.executeParallelJobSearch(functionCalls), nil
}

func (a *AIClient) executeParallelJobSearch(functionCalls []*genai.FunctionCall) []dtos.Job {
	var wg sync.WaitGroup
	results := make(chan dtos.JobSearchResult, len(functionCalls))

	for _, functionCall := range functionCalls {
		wg.Add(1)
		go func(fc *genai.FunctionCall) {
			defer wg.Done()

			query, ok := fc.Args["query"].(string)
			if !ok {
				fmt.Printf("No query found for function %s\n", fc.Name)
				results <- dtos.JobSearchResult{
					Jobs:   []dtos.Job{},
					Error:  fmt.Errorf("no query found for function %s", fc.Name),
					Source: fc.Name,
				}
				return
			}

			fmt.Printf("Searching %s for: %s\n", fc.Name, query)

			switch fc.Name {
			case "search_jsearch_jobs":
				jobs, err := SearchJobsJSearch(query)
				results <- dtos.JobSearchResult{
					Jobs:   jobs,
					Error:  err,
					Source: "JSearch",
				}

			case "search_structured_jobs":
				structuredJobs, err := SearchJobsLinkUpStructured(query)
				if err != nil {
					results <- dtos.JobSearchResult{
						Jobs:   []dtos.Job{},
						Error:  err,
						Source: "LinkUp-Structured",
					}
				} else {
					jobs := convertStructuredToRegularJobs(structuredJobs)
					results <- dtos.JobSearchResult{
						Jobs:   jobs,
						Error:  nil,
						Source: "LinkUp-Structured",
					}
				}

			default:
				fmt.Printf("Unknown function: %s\n", fc.Name)
				results <- dtos.JobSearchResult{
					Jobs:   []dtos.Job{},
					Error:  fmt.Errorf("unknown function: %s", fc.Name),
					Source: fc.Name,
				}
			}
		}(functionCall)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allJobs []dtos.Job
	for result := range results {
		if result.Error != nil {
			fmt.Printf("%s error: %v\n", result.Source, result.Error)
		} else {
			fmt.Printf("%s found %d jobs\n", result.Source, len(result.Jobs))
			allJobs = append(allJobs, result.Jobs...)
		}
	}

	return allJobs
}

func (a *AIClient) parseLocationPreference(locationPreference dtos.LocationPreference) string {
	locationContext := ""
	if len(locationPreference.Types) > 0 {
		locationContext = fmt.Sprintf("\nWork arrangement preferences: %s", strings.Join(locationPreference.Types, ", "))
		if len(locationPreference.Locations) > 0 {
			locationContext += fmt.Sprintf("\nPreferred locations: %s", strings.Join(locationPreference.Locations, ", "))
		}
	}

	return locationContext
}

func convertStructuredToRegularJobs(structuredJobs *dtos.JobAnnouncements) []dtos.Job {
	var jobs []dtos.Job

	for _, structJob := range structuredJobs.Jobs {
		location := "Remote"
		if !structJob.Remote && structJob.Location != nil {
			location = *structJob.Location
		}

		description := fmt.Sprintf("Experience Level: %s\nRequired Skills: %s",
			structJob.ExperienceLevel,
			strings.Join(structJob.RequiredSkills, ", "),
		)

		if structJob.Salary != nil {
			description += fmt.Sprintf("\nSalary: $%d", *structJob.Salary)
		}

		jobs = append(jobs, dtos.Job{
			Title:       structJob.JobTitle,
			Company:     structJob.Company,
			Location:    location,
			Description: description,
			URL:         structJob.JobPostURL,
			Source:      "LinkUp-Structured",
		})
	}

	return jobs
}
