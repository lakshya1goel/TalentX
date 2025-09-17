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

	candidateProfile, err := r.extractCandidateProfile(ctx, pdfBytes)
	if err != nil {
		fmt.Printf("Error extracting candidate profile: %v\n", err)
		return r.fallbackRanking(jobs), nil
	}

	fmt.Println("Candidate profile: ", candidateProfile)

	if len(jobs) > 10 {
		return r.RerankJobsParallel(ctx, candidateProfile, jobs)
	}

	return r.rankBatchJobs(ctx, candidateProfile, jobs)
}

func (r *RerankingClient) extractCandidateProfile(ctx context.Context, pdfBytes []byte) (string, error) {
	prompt := `
You are a resume analyzer. Extract and summarize the candidate's profile from their resume.

Provide a comprehensive summary including:
- Job titles they would be suitable for
- Technical skills and expertise areas
- Years of professional experience
- Education background
- Location preferences (if mentioned)
- Work location preferences (remote/hybrid/on-site)
- Industry experience
- Notable projects or achievements

NOTE: Do not consider internships as professional experience. They are treated as fresher. So look for internships and entry level jobs.

Format this as a clear, structured profile summary.
`

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
		"gemini-1.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to extract candidate profile: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no profile response from AI")
	}

	profileText := ""
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			profileText += part.Text
		}
	}

	return profileText, nil
}

func (r *RerankingClient) RerankJobsParallel(ctx context.Context, candidateProfile string, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	const batchSize = 10
	const maxConcurrency = 3

	var batches [][]dtos.Job
	for i := 0; i < len(jobs); i += batchSize {
		end := min(i+batchSize, len(jobs))
		batches = append(batches, jobs[i:end])
	}

	monitorConcurrency := make(chan struct{}, maxConcurrency)
	results := make(chan dtos.BatchResult, len(batches))
	var wg sync.WaitGroup

	for batchIndex, batch := range batches {
		wg.Add(1)
		go func(idx int, jobBatch []dtos.Job) {
			defer wg.Done()

			monitorConcurrency <- struct{}{}
			defer func() { <-monitorConcurrency }()

			fmt.Printf("Processing batch %d with %d jobs\n", idx+1, len(jobBatch))

			rankedBatch, err := r.rankBatchJobs(ctx, candidateProfile, jobBatch)
			if err != nil {
				fmt.Printf("Error ranking batch %d: %v, using fallback\n", idx+1, err)
				rankedBatch = r.fallbackRanking(jobBatch)
			}

			fmt.Printf("Completed batch %d with %d ranked jobs\n", idx+1, len(rankedBatch))

			results <- dtos.BatchResult{
				Jobs:       rankedBatch,
				BatchIndex: idx,
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

func (r *RerankingClient) rankBatchJobs(ctx context.Context, candidateProfile string, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	if len(jobs) == 0 {
		return []dtos.RankedJob{}, nil
	}

	var jobsJSON []string
	for i, job := range jobs {
		jobJSON, err := json.MarshalIndent(job, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling job %d: %v\n", i, err)
			continue
		}
		jobsJSON = append(jobsJSON, string(jobJSON))
	}

	systemMessage := `You are a job matching assistant. Your task is to evaluate multiple jobs based on their match with the candidate's profile, taking into account the job title, the skills required, the seniority level, the physical location (where the company offering the work is based in) and the working location (remote/hybrid/on-site). You then have to produce a match score (between 0 and 100) for each job and justify that match score explaining your reasons.

	Evaluation Criteria:
	1. Job title alignment with candidate's potential roles
	2. Required skills match with candidate's skills
	3. Seniority level alignment (internship, entry level, junior, mid-level, senior)
	4. Location preferences and work arrangement compatibility
	5. Industry and domain experience relevance
	6. Overall career trajectory fit

	Provide your evaluation in the following JSON format for ALL jobs:
	{
		"evaluations": [
			{
				"job_index": 0,
				"match_score": <integer between 0-100>,
				"reasons": "<detailed explanation of the match evaluation>",
				"skills_matched": ["<list of matched skills>"],
				"experience_match": "<assessment of experience level fit>"
			},
			{
				"job_index": 1,
				"match_score": <integer between 0-100>,
				"reasons": "<detailed explanation of the match evaluation>",
				"skills_matched": ["<list of matched skills>"],
				"experience_match": "<assessment of experience level fit>"
			}
		]
	}

	IMPORTANT: Provide evaluations for ALL jobs in the same order they are presented. Use job_index to match each evaluation to its corresponding job.`

	jobsListStr := ""
	for i, jobJSON := range jobsJSON {
		jobsListStr += fmt.Sprintf("Job %d:\n%s\n\n", i, jobJSON)
	}

	userMessage := fmt.Sprintf(`System: %s
		User: Here is my profile:
		'''
		%s
		'''
		And here are the JSON cards of %d jobs that I found:
		'''
		%s
		'''
		Can you evaluate the match for all these jobs?`, systemMessage, candidateProfile, len(jobs), jobsListStr)

	parts := []*genai.Part{
		genai.NewPartFromText(userMessage),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	temp := float32(0.1)
	result, err := r.Client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature: &temp,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate job batch: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no evaluation response from AI")
	}

	responseText := ""
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	rankedJobs, err := r.parseBatchJobEvaluation(responseText, jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse batch evaluation response: %w", err)
	}

	var filteredJobs []dtos.RankedJob
	for _, rankedJob := range rankedJobs {
		if rankedJob.PercentMatch >= 30.0 {
			filteredJobs = append(filteredJobs, rankedJob)
		}
	}

	sort.Slice(filteredJobs, func(i, j int) bool {
		return filteredJobs[i].PercentMatch > filteredJobs[j].PercentMatch
	})

	return filteredJobs, nil
}

func (r *RerankingClient) parseBatchJobEvaluation(responseText string, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}") + 1

	if jsonStart == -1 || jsonEnd == 0 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := responseText[jsonStart:jsonEnd]

	var batchEvaluation struct {
		Evaluations []struct {
			JobIndex        int      `json:"job_index"`
			MatchScore      int      `json:"match_score"`
			Reasons         string   `json:"reasons"`
			SkillsMatched   []string `json:"skills_matched"`
			ExperienceMatch string   `json:"experience_match"`
		} `json:"evaluations"`
	}

	err := json.Unmarshal([]byte(jsonStr), &batchEvaluation)
	if err != nil {
		return nil, fmt.Errorf("error parsing batch evaluation JSON: %w", err)
	}

	var rankedJobs []dtos.RankedJob
	for _, evaluation := range batchEvaluation.Evaluations {
		if evaluation.JobIndex < 0 || evaluation.JobIndex >= len(jobs) {
			fmt.Printf("Invalid job index %d, skipping\n", evaluation.JobIndex)
			continue
		}

		rankedJob := dtos.RankedJob{
			Job:             jobs[evaluation.JobIndex],
			PercentMatch:    float64(evaluation.MatchScore),
			MatchReason:     evaluation.Reasons,
			SkillsMatched:   evaluation.SkillsMatched,
			ExperienceMatch: evaluation.ExperienceMatch,
		}
		rankedJobs = append(rankedJobs, rankedJob)
	}

	if len(rankedJobs) < len(jobs) {
		fmt.Printf("Warning: Only got %d evaluations for %d jobs, using fallback for missing ones\n", len(rankedJobs), len(jobs))

		evaluatedJobs := make(map[int]bool)
		for _, evaluation := range batchEvaluation.Evaluations {
			if evaluation.JobIndex >= 0 && evaluation.JobIndex < len(jobs) {
				evaluatedJobs[evaluation.JobIndex] = true
			}
		}

		for i := 0; i < len(jobs); i++ {
			if !evaluatedJobs[i] {
				rankedJob := dtos.RankedJob{
					Job:             jobs[i],
					PercentMatch:    50.0,
					MatchReason:     "Fallback evaluation - AI did not provide evaluation for this job",
					SkillsMatched:   []string{},
					ExperienceMatch: "Unknown",
				}
				rankedJobs = append(rankedJobs, rankedJob)
			}
		}
	}

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
