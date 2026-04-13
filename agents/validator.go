package agents

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/blackviking27/validate-idea-cli/providers"
	"github.com/blackviking27/validate-idea-cli/tools"
)

type ValidationResult struct {
	EnhancedIdea string
	AuditReport  string
	GrowthReport string
}

func renderTemplate(path string, data interface{}) (string, error) {
	templateContent, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	template, err := template.New("tmpelate").Parse(string(templateContent))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err = template.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func generateAIReadyPostResult(results []tools.ParsedSearchResult) string {
	var parsedResult strings.Builder
	for _, res := range results {

		comments := ""
		for _, cmt := range res.Comments {
			comments += fmt.Sprintf("%s,", cmt)
		}

		post := fmt.Sprintf("Title: %s\nContent: %s\nComments:%s\n\n", res.Title, res.Content, comments)
		parsedResult.WriteString(post)
	}

	return parsedResult.String()
}

func RunValidator(ctx context.Context, aiModel *providers.AIProvider, userIdea string) (ValidationResult, error) {
	validationResult := ValidationResult{}

	// --- PHASE 1: Summary ---
	// Generating user idea
	fmt.Printf("🔍 Analyzing idea with %s...\n\n", (*aiModel).Name())
	enhanceIdeaPrompt, err := renderTemplate("prompts/generate-idea.txt", struct{ Idea string }{Idea: userIdea})
	if err != nil {
		log.Fatalf("Failed to load generate idea prompt, err: %v", err)
	}

	enhancedIdea, err := (*aiModel).Generate(ctx, enhanceIdeaPrompt)
	if err != nil {
		return validationResult, err
	}
	validationResult.EnhancedIdea = enhancedIdea

	// --- PHASE 2: Reddit Search & Audit ---
	// Running the platform validation
	fmt.Println("\n\n🌐 Searching Reddit for community feedback...")
	searchQuery, err := (*aiModel).Generate(ctx, fmt.Sprintf("Generate a short 3 word search query for Reddit to find if people need: %s", enhancedIdea))
	if err != nil || searchQuery == "" {
		return validationResult, fmt.Errorf("Unable to generate search query, err: %v", err)
	}

	redditSearch := tools.NewRedditSearch()
	searchResults, err := redditSearch.Search(ctx, searchQuery)
	if err != nil {
		return validationResult, fmt.Errorf("Search failed for query: %v", searchQuery)
	}

	aiReadyPostResult := generateAIReadyPostResult(searchResults)

	// Generating audit report
	fmt.Println("\n\nGenerating audit report...")
	auditPrompt, err := renderTemplate("prompts/research-prompt.txt", struct {
		Idea    string
		Results string
	}{Idea: userIdea, Results: aiReadyPostResult})

	auditResult, err := (*aiModel).Generate(ctx, auditPrompt)
	if err != nil {
		return validationResult, fmt.Errorf("Unable to generate audit report, err: %v", err)
	}
	validationResult.AuditReport = auditResult

	// --- PHASE 3 & 4: Keywords & Discovery ---
	fmt.Println("📈 Generating discovery strategy & keywords...")
	growthPrompt, _ := renderTemplate("prompts/growth_stratergy.txt", struct {
		Summary string
		Audit   string
	}{Summary: enhancedIdea, Audit: auditResult})
	growthReport, err := (*aiModel).Generate(ctx, growthPrompt)
	if err != nil {
		fmt.Println("Unable to generate growth report", err)
		return validationResult, err
	}
	validationResult.GrowthReport = growthReport

	return validationResult, nil
}
