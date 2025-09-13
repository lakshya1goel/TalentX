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
	prompt := a.PromptWithLocation(locationPreference)

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

	fmt.Println("Processing AI response...")
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

// New method: Parse resume and get structured jobs from LinkUp
func (a *AIClient) GetStructuredJobsFromResume(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.Job, error) {
	fmt.Println("Parsing resume to generate job description...")

	// Step 1: Parse resume and generate job description
	jobDescription, err := a.generateJobDescriptionFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return nil, fmt.Errorf("failed to generate job description from resume: %w", err)
	}

	fmt.Printf("Generated job description: %s\n", jobDescription)

	// Step 2: Search LinkUp with structured output using the generated description
	fmt.Println("Searching LinkUp with structured output...")
	structuredJobs, err := SearchJobsLinkUpStructured(jobDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to search structured jobs: %w", err)
	}

	// Step 3: Convert structured jobs to regular jobs for compatibility with existing ranking
	jobs := convertStructuredToRegularJobs(structuredJobs)

	fmt.Printf("Found %d structured jobs from LinkUp\n", len(jobs))
	return jobs, nil
}

// Generate job description from resume for LinkUp search
func (a *AIClient) generateJobDescriptionFromResume(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) (string, error) {
	locationContext := ""
	if len(locationPreference.Types) > 0 {
		locationContext = fmt.Sprintf("\nWork arrangement preferences: %s", strings.Join(locationPreference.Types, ", "))
		if len(locationPreference.Locations) > 0 {
			locationContext += fmt.Sprintf("\nPreferred locations: %s", strings.Join(locationPreference.Locations, ", "))
		}
	}

	prompt := fmt.Sprintf(`You are a professional resume analyzer and job search expert. Analyze the provided resume and generate a comprehensive job search description that will be used to find relevant job opportunities.

Your task is to:
1. Extract the candidate's key skills, experience level, and expertise areas
2. Identify their primary job roles and career focus
3. Generate a detailed job search query that captures what they're looking for

Generate a job search description that includes:
- Primary role/job title they would be suitable for
- Key technical skills and technologies they know
- Experience level (entry-level, junior, mid-level, senior, lead)
- Industry preferences if evident from their background
- Any specializations or domain expertise

%s

Format your response as a single, comprehensive job search description that would be effective for finding relevant positions. Focus on being specific about skills and experience level while being broad enough to capture multiple relevant opportunities.

Example format: "Senior Software Engineer with 5+ years experience in Python, Django, React, and AWS. Looking for backend or full-stack roles in fintech or healthcare. Strong experience with microservices, API development, and cloud infrastructure."

Return only the job search description, no additional text.`, locationContext)

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

	temp := float32(0.3)
	result, err := a.Client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate job description: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from AI for job description generation")
	}

	jobDescription := ""
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			jobDescription += part.Text
		}
	}

	if jobDescription == "" {
		return "", fmt.Errorf("empty job description generated")
	}

	return strings.TrimSpace(jobDescription), nil
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
					// Convert structured jobs to regular jobs for compatibility
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

	fmt.Printf("Total jobs found across all sources: %d\n", len(allJobs))
	return allJobs
}

// Helper function to convert structured jobs to regular jobs for compatibility
func convertStructuredToRegularJobs(structuredJobs *dtos.JobAnnouncements) []dtos.Job {
	var jobs []dtos.Job

	for _, structJob := range structuredJobs.Jobs {
		location := "Remote"
		if !structJob.Remote && structJob.Location != nil {
			location = *structJob.Location
		}

		// Create a description from structured data
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
