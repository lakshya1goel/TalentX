package service

import (
	"context"
	"fmt"

	"github.com/lakshya1goel/job-assistance/config"
	"github.com/lakshya1goel/job-assistance/internal/ai"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobService interface {
	FetchAndRankAllJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.RankedJob, error)
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

func (s *jobService) FetchAndRankAllJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.RankedJob, error) {
	fmt.Println("Fetching jobs from various sources...")
	jobs, err := s.aiClient.GetJobsFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}

	if len(jobs) == 0 {
		return []dtos.RankedJob{}, nil
	}

	fmt.Printf("Ranking %d jobs based on resume relevance...\n", len(jobs))
	rankedJobs, err := s.rankingClient.RerankJobs(ctx, pdfBytes, jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to rank jobs: %w", err)
	}

	fmt.Printf("Successfully processed and ranked %d jobs\n", len(rankedJobs))
	return rankedJobs, nil
}
