package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
	"google.golang.org/genai"
)

type RerankingClient struct {
	Client *genai.Client
}

func NewRerankingClient(ctx context.Context, apiKey string) *RerankingClient {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		fmt.Printf("Error creating reranking client: %v\n", err)
	}

	return &RerankingClient{Client: client}
}

func (r *RerankingClient) RerankJobs(ctx context.Context, pdfBytes []byte, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	if len(jobs) == 0 {
		return []dtos.RankedJob{}, nil
	}

	if len(jobs) > 60 {
		fmt.Printf("Limiting ranking to first 60 jobs out of %d total jobs\n", len(jobs))
		jobs = jobs[:60]
	}

	if len(jobs) > 10 {
		return r.RerankJobsParallel(ctx, pdfBytes, jobs)
	}

	return r.rerankJobsSingle(ctx, pdfBytes, jobs)
}

func (r *RerankingClient) RerankJobsParallel(ctx context.Context, pdfBytes []byte, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	const batchSize = 10
	const maxConcurrency = 3

	var batches [][]dtos.Job
	for i := 0; i < len(jobs); i += batchSize {
		end := i + batchSize
		if end > len(jobs) {
			end = len(jobs)
		}
		batches = append(batches, jobs[i:end])
	}

	semaphore := make(chan struct{}, maxConcurrency)
	results := make(chan dtos.BatchResult, len(batches))
	var wg sync.WaitGroup

	for batchIndex, batch := range batches {
		wg.Add(1)
		go func(batchIdx int, jobBatch []dtos.Job) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fmt.Printf("Processing batch %d with %d jobs\n", batchIdx+1, len(jobBatch))

			rankedBatch, err := r.rerankJobsSingle(ctx, pdfBytes, jobBatch)
			if err != nil {
				fmt.Printf("Error ranking batch %d: %v, using fallback\n", batchIdx+1, err)
				rankedBatch = r.fallbackRanking(jobBatch)
			}

			fmt.Printf("Completed batch %d with %d ranked jobs\n", batchIdx+1, len(rankedBatch))

			results <- dtos.BatchResult{
				Jobs:       rankedBatch,
				BatchIndex: batchIdx,
			}
		}(batchIndex, batch)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allRanked []dtos.RankedJob
	batchCount := 0
	for result := range results {
		allRanked = append(allRanked, result.Jobs...)
		batchCount++
		fmt.Printf("Collected results from batch %d (%d/%d batches complete)\n", result.BatchIndex+1, batchCount, len(batches))
	}

	sort.Slice(allRanked, func(i, j int) bool {
		return allRanked[i].PercentMatch > allRanked[j].PercentMatch
	})

	fmt.Printf("Successfully ranked and sorted %d jobs from %d batches\n", len(allRanked), len(batches))
	return allRanked, nil
}

func (r *RerankingClient) rerankJobsSingle(ctx context.Context, pdfBytes []byte, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	prompt := r.RerankingPrompt(jobs)

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

	temp := float32(0.1)
	result, err := r.Client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to rerank jobs: %w", err)
	}

	return r.parseRankingResponse(result, jobs)
}

func (r *RerankingClient) parseRankingResponse(result *genai.GenerateContentResponse, originalJobs []dtos.Job) ([]dtos.RankedJob, error) {
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	responseText := ""
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	jsonStart := strings.Index(responseText, "[")
	jsonEnd := strings.LastIndex(responseText, "]") + 1

	if jsonStart == -1 || jsonEnd == 0 {
		return r.fallbackRanking(originalJobs), nil
	}

	jsonStr := responseText[jsonStart:jsonEnd]

	var rankings []struct {
		JobIndex        int      `json:"job_index"`
		MatchScore      float64  `json:"match_score"`
		MatchReason     string   `json:"match_reason"`
		SkillsMatched   []string `json:"skills_matched"`
		ExperienceMatch string   `json:"experience_match"`
		Concerns        string   `json:"concerns,omitempty"`
	}

	err := json.Unmarshal([]byte(jsonStr), &rankings)
	if err != nil {
		fmt.Printf("Error parsing ranking JSON: %v\n", err)
		return r.fallbackRanking(originalJobs), nil
	}

	var rankedJobs []dtos.RankedJob
	for _, ranking := range rankings {
		percentMatch := ranking.MatchScore * 100.0

		if ranking.JobIndex >= 0 && ranking.JobIndex < len(originalJobs) && percentMatch >= 30.0 {
			matchReason := ranking.MatchReason
			if ranking.Concerns != "" {
				matchReason += " | Concerns: " + ranking.Concerns
			}

			rankedJobs = append(rankedJobs, dtos.RankedJob{
				Job:             originalJobs[ranking.JobIndex],
				PercentMatch:    percentMatch,
				MatchReason:     matchReason,
				SkillsMatched:   ranking.SkillsMatched,
				ExperienceMatch: ranking.ExperienceMatch,
			})
		}
	}

	fmt.Printf("Successfully ranked %d jobs\n", len(rankedJobs))
	return rankedJobs, nil
}

func (r *RerankingClient) fallbackRanking(jobs []dtos.Job) []dtos.RankedJob {
	fmt.Println("Using fallback ranking - returning jobs in original order")
	var rankedJobs []dtos.RankedJob

	for i, job := range jobs {
		percentage := 80.0 - float64(i)*5.0

		if percentage >= 30.0 {
			rankedJobs = append(rankedJobs, dtos.RankedJob{
				Job:             job,
				PercentMatch:    percentage,
				MatchReason:     "Fallback ranking - AI parsing failed",
				SkillsMatched:   []string{},
				ExperienceMatch: "Unknown",
			})
		}
	}

	return rankedJobs
}
