package service

import (
	"context"
	"fmt"

	"github.com/lakshya1goel/job-assistance/internal/ai"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobService interface {
	FetchAndRankStructuredJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference, apiKey string) ([]dtos.RankedJob, error)
}

type jobService struct {
}

func NewJobService() JobService {
	return &jobService{}
}

func (s *jobService) FetchAndRankStructuredJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference, apiKey string) ([]dtos.RankedJob, error) {
	fmt.Println("Parsing resume and searching with structured output...")

	aiClient := ai.NewAIClient(ctx, apiKey)
	rankingClient := ai.NewRerankingClient(ctx, apiKey)

	jobs, err := aiClient.GetJobsFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return nil, fmt.Errorf("failed to get structured jobs from resume: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No jobs found from structured search")
		return []dtos.RankedJob{}, nil
	}

	fmt.Printf("Re-ranking %d structured jobs based on resume relevance...\n", len(jobs))
	rankedJobs, err := rankingClient.RerankJobs(ctx, pdfBytes, jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to rank structured jobs: %w", err)
	}
	return rankedJobs, nil
}
