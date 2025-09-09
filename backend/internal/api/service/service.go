package service

import (
	"context"
	"fmt"
	"math"

	"github.com/lakshya1goel/job-assistance/config"
	"github.com/lakshya1goel/job-assistance/internal/ai"
	"github.com/lakshya1goel/job-assistance/internal/dtos"
)

type JobService interface {
	FetchJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference, pagination dtos.PaginationRequest) (dtos.PaginatedJobResponse, error)
}

type jobService struct {
	aiClient      ai.AIClient
	rankingClient ai.RerankingClient
}

func NewJobService() JobService {
	return &jobService{
		aiClient:      *ai.NewAIClient(context.Background(), config.GetAPIKey()),
		rankingClient: *ai.NewRerankingClient(context.Background(), config.GetAPIKey()),
	}
}

func (s *jobService) FetchJobs(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference, pagination dtos.PaginationRequest) (dtos.PaginatedJobResponse, error) {
	fmt.Println("Fetching jobs from various sources...")
	jobs, err := s.aiClient.GetJobsFromResume(ctx, pdfBytes, locationPreference)
	if err != nil {
		return dtos.PaginatedJobResponse{}, fmt.Errorf("failed to analyze resume: %w", err)
	}

	if len(jobs) == 0 {
		return dtos.PaginatedJobResponse{
			Jobs:       []dtos.RankedJob{},
			TotalJobs:  0,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: 0,
			Success:    true,
		}, nil
	}

	fmt.Printf("Ranking %d jobs based on resume relevance...\n", len(jobs))
	rankedJobs, err := s.rankingClient.RerankJobs(ctx, pdfBytes, jobs)
	if err != nil {
		return dtos.PaginatedJobResponse{}, fmt.Errorf("failed to rank jobs: %w", err)
	}

	totalJobs := len(rankedJobs)
	totalPages := int(math.Ceil(float64(totalJobs) / float64(pagination.PageSize)))

	startIndex := (pagination.Page - 1) * pagination.PageSize
	endIndex := startIndex + pagination.PageSize

	if startIndex >= totalJobs {
		return dtos.PaginatedJobResponse{
			Jobs:       []dtos.RankedJob{},
			TotalJobs:  totalJobs,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: totalPages,
			Success:    true,
		}, nil
	}

	if endIndex > totalJobs {
		endIndex = totalJobs
	}

	paginatedJobs := rankedJobs[startIndex:endIndex]
	
	fmt.Printf("Returning page %d of %d (jobs %d-%d of %d total)\n", 
		pagination.Page, totalPages, startIndex+1, endIndex, totalJobs)

	return dtos.PaginatedJobResponse{
		Jobs:       paginatedJobs,
		TotalJobs:  totalJobs,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
		Success:    true,
	}, nil
}
