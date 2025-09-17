package ai

import (
	"context"
	"fmt"

	"github.com/lakshya1goel/job-assistance/internal/dtos"
	"google.golang.org/genai"
)

type ProfileClient struct {
	Client *genai.Client
}

func NewProfileClient(ctx context.Context, apiKey string) *ProfileClient {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		fmt.Printf("Error creating profile client: %v\n", err)
	}

	return &ProfileClient{
		Client: client,
	}
}

func (p *ProfileClient) ExtractCandidateProfile(ctx context.Context, pdfBytes []byte, locationPreference dtos.LocationPreference) (string, error) {
	prompt := p.CandidateProfilePrompt(locationPreference)

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
	result, err := p.Client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
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

	if profileText == "" {
		return "", fmt.Errorf("empty profile response from AI")
	}

	return profileText, nil
}
