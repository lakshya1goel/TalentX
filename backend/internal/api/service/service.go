package service

import (
	"context"
	"fmt"

	"github.com/lakshya1goel/job-assistance/config"
	"github.com/lakshya1goel/job-assistance/internal/ai"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobService interface {
	FetchJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.Job, error)
}

type jobService struct {
	aiClient ai.AIClient
}

func NewJobService() JobService {
	return &jobService{
		aiClient: *ai.NewAIClient(context.Background(), config.GetAPIKey()),
	}
}

func (s *jobService) FetchJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) ([]dtos.Job, error) {
	jobs, err := s.aiClient.GetJobsFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return []dtos.Job{}, fmt.Errorf("failed to analyze resume: %w", err)
	}

	return jobs, nil
}
