package ai

import (
	"context"
	"fmt"

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

	var allJobs []dtos.Job

	fmt.Println("Processing AI response...")

	for _, candidate := range result.Candidates {
		for _, part := range candidate.Content.Parts {

			if part.FunctionCall != nil {
				functionCall := part.FunctionCall
				fmt.Printf("AI wants to call: %s\n", functionCall.Name)

				jobs := a.callJobAPIWithLocation(functionCall)
				allJobs = append(allJobs, jobs...)
			}
		}
	}

	return allJobs, nil
}

func (a *AIClient) callJobAPIWithLocation(functionCall *genai.FunctionCall) []dtos.Job {
	query, ok := functionCall.Args["query"].(string)
	if !ok {
		fmt.Println("No query found")
		return []dtos.Job{}
	}

	fmt.Printf("Searching for: %s\n", query)

	switch functionCall.Name {
	case "search_jsearch_jobs":
		jobs, err := SearchJobsJSearch(query)
		if err != nil {
			fmt.Printf("JSearch error: %v\n", err)
			return []dtos.Job{}
		}
		fmt.Printf("JSearch found %d jobs\n", len(jobs))
		return jobs

	case "search_linkup_jobs":
		jobs, err := SearchJobsLinkUp(query)
		if err != nil {
			fmt.Printf("LinkUp error: %v\n", err)
			return []dtos.Job{}
		}
		fmt.Printf("LinkUp found %d jobs\n", len(jobs))
		return jobs

	default:
		fmt.Printf("Unknown function: %s\n", functionCall.Name)
		return []dtos.Job{}
	}
}
