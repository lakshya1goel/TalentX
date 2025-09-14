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

	if len(jobs) > 10 {
		return r.RerankJobsParallel(ctx, candidateProfile, jobs)
	}

	return r.rerankJobsSingle(ctx, candidateProfile, jobs)
}

func (r *RerankingClient) extractCandidateProfile(ctx context.Context, pdfBytes []byte) (string, error) {
	prompt := `
You are a resume analyzer. Extract and summarize the candidate's profile from their resume.

Provide a comprehensive summary including:
- Job titles they would be suitable for
- Seniority level (internship, entry level, junior, mid-level, senior)
- Technical skills and expertise areas
- Years of professional experience
- Education background
- Location preferences (if mentioned)
- Work location preferences (remote/hybrid/on-site)
- Industry experience
- Notable projects or achievements

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

			rankedBatch, err := r.rerankJobsSingle(ctx, candidateProfile, jobBatch)
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

func (r *RerankingClient) rerankJobsSingle(ctx context.Context, candidateProfile string, jobs []dtos.Job) ([]dtos.RankedJob, error) {
	var rankedJobs []dtos.RankedJob

	for i, job := range jobs {
		evaluation, err := r.evaluateJobMatch(ctx, candidateProfile, job)
		if err != nil {
			fmt.Printf("Error evaluating job %d: %v\n", i, err)
			continue
		}

		if evaluation != nil && evaluation.PercentMatch >= 30.0 {
			rankedJobs = append(rankedJobs, *evaluation)
		}
	}

	sort.Slice(rankedJobs, func(i, j int) bool {
		return rankedJobs[i].PercentMatch > rankedJobs[j].PercentMatch
	})

	return rankedJobs, nil
}

func (r *RerankingClient) evaluateJobMatch(ctx context.Context, candidateProfile string, job dtos.Job) (*dtos.RankedJob, error) {
	systemMessage := `You are a job matching assistant. Your task is to evaluate a job based on its match with the candidate's profile, taking into account the job title, the skills required, the seniority level, the physical location (where the company offering the work is based in) and the working location (remote/hybrid/on-site). You then have to produce a match score (between 0 and 100) and justify that match score explaining your reasons for that.

	Evaluation Criteria:
	1. Job title alignment with candidate's potential roles
	2. Required skills match with candidate's skills
	3. Seniority level alignment (internship, entry level, junior, mid-level, senior)
	4. Location preferences and work arrangement compatibility
	5. Industry and domain experience relevance
	6. Overall career trajectory fit

	Provide your evaluation in the following JSON format:
	{
		"match_score": <integer between 0-100>,
		"reasons": "<detailed explanation of the match evaluation>",
		"skills_matched": ["<list of matched skills>"],
		"experience_match": "<assessment of experience level fit>"
	}`

		jobJSON, _ := json.MarshalIndent(job, "", "  ")

		userMessage := fmt.Sprintf(`System: %s

	User: Here is my profile:

	'''
	%s
	'''

	And here is the JSON card of a job that I found:

	'''
	%s
	'''

	Can you evaluate the match for me?`, systemMessage, candidateProfile, string(jobJSON))

	parts := []*genai.Part{
		genai.NewPartFromText(userMessage),
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
		return nil, fmt.Errorf("failed to evaluate job match: %w", err)
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

	evaluation, err := r.parseJobEvaluation(responseText, job)
	if err != nil {
		return nil, fmt.Errorf("failed to parse evaluation response: %w", err)
	}

	return evaluation, nil
}

func (r *RerankingClient) parseJobEvaluation(responseText string, job dtos.Job) (*dtos.RankedJob, error) {
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}") + 1

	if jsonStart == -1 || jsonEnd == 0 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := responseText[jsonStart:jsonEnd]

	var evaluation struct {
		MatchScore      int      `json:"match_score"`
		Reasons         string   `json:"reasons"`
		SkillsMatched   []string `json:"skills_matched"`
		ExperienceMatch string   `json:"experience_match"`
	}

	err := json.Unmarshal([]byte(jsonStr), &evaluation)
	if err != nil {
		return nil, fmt.Errorf("error parsing evaluation JSON: %w", err)
	}

	return &dtos.RankedJob{
		Job:             job,
		PercentMatch:    float64(evaluation.MatchScore),
		MatchReason:     evaluation.Reasons,
		SkillsMatched:   evaluation.SkillsMatched,
		ExperienceMatch: evaluation.ExperienceMatch,
	}, nil
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
