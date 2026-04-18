package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type GeminiProvider struct {
	client *googleai.GoogleAI
}

func (this *GeminiProvider) Name() string {
	return "Gemini"
}

func (this *GeminiProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, this.client, prompt)
}

func NewGeminiProvider(ctx context.Context, model string) (*GeminiProvider, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not found in environment")
	}

	client, err := googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(model))
	if err != nil {
		return nil, err
	}

	return &GeminiProvider{
		client: client,
	}, nil
}
