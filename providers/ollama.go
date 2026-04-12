package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaProvider struct {
	client *ollama.LLM
}

func (this *OllamaProvider) Name() string { return "ollama" }

func (this *OllamaProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, this.client, prompt)
}

func NewOllamaProvider(ctx context.Context) (*OllamaProvider, error) {
	ollamaHost := os.Getenv("OLLAMA_HOST")

	if ollamaHost == "" {
		return nil, fmt.Errorf("No ollama host defined in environment")
	}

	client, err := ollama.New(ollama.WithServerURL(ollamaHost), ollama.WithModel("gemma4:e2b"))
	if err != nil {
		return nil, err
	}

	return &OllamaProvider{
		client: client,
	}, nil

}
