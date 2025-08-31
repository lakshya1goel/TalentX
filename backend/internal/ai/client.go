package ai

import (
	"context"
	"fmt"

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
