package agents

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/blackviking27/validate-idea-cli/providers"
)

func RunValidator(ctx context.Context, aiModel *providers.AIProvider, userIdea string) (string, error) {
	// Loading the prompt
	validateIdeaPrompt, err := os.ReadFile("prompts/validate-idea.txt")
	if err != nil {
		return "", nil
	}

	tmpl, err := template.New("ValidateIdea").Parse(string(validateIdeaPrompt))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	data := struct{ Idea string }{Idea: userIdea}
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	return (*aiModel).Generate(ctx, buffer.String())
}
