package service

import (
	"context"
	"fmt"

	"github.com/lakshya1goel/job-assistance/config"
	"github.com/lakshya1goel/job-assistance/internal/ai"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobService interface {
	FetchAndRankStructuredJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.RankedJob, error)
}

type jobService struct {
	aiClient      *ai.AIClient
	rankingClient *ai.RerankingClient
}

func NewJobService() JobService {
	ctx := context.Background()
	apiKey := config.GetAPIKey()

	return &jobService{
		aiClient:      ai.NewAIClient(ctx, apiKey),
		rankingClient: ai.NewRerankingClient(ctx, apiKey),
	}
}

func (s *jobService) FetchAndRankStructuredJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.RankedJob, error) {
	fmt.Println("Parsing resume and searching with structured output...")
	jobs, err := s.aiClient.GetJobsFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return nil, fmt.Errorf("failed to get structured jobs from resume: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No jobs found from structured search")
		return []dtos.RankedJob{}, nil
	}

	fmt.Printf("Re-ranking %d structured jobs based on resume relevance...\n", len(jobs))
	rankedJobs, err := s.rankingClient.RerankJobs(ctx, pdfBytes, jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to rank structured jobs: %w", err)
	}
	return rankedJobs, nil
}
