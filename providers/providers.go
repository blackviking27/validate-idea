package providers

import "context"

type AIProvider interface {
	Name() string                                                // Return the provider name
	Generate(ctx context.Context, prompt string) (string, error) // Handles simple text-to-text prompt
}
