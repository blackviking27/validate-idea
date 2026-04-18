package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blackviking27/validate-idea-cli/agents"
	"github.com/blackviking27/validate-idea-cli/providers"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	providerArg string
	aiModel     string
)

func getAiProvider(ctx context.Context, config ValidationConfig) *providers.AIProvider {
	var currProivder providers.AIProvider

	switch config.Provider {
	case "Gemini":
	case "gemini":
		currProivder, _ = providers.NewGeminiProvider(ctx, config.Model)
	case "ollama":
	case "Ollama":
		currProivder, _ = providers.NewOllamaProvider(ctx, config.Model)
	}

	if currProivder == nil {
		log.Fatal("Not a valid AI provider provided: ", config.Provider)
	}

	return &currProivder
}

func runValidateCommand(cmd *cobra.Command, config ValidationConfig) agents.ValidationResult {
	ctx := context.Background()
	userIdea := config.Idea

	// Initiating ai provider
	aiProvider := getAiProvider(ctx, config)

	result, err := agents.RunValidator(ctx, aiProvider, userIdea)
	if err != nil {
		log.Fatalf("Validation failed : %v", err)
	}

	return result
}

func saveReport(content agents.ValidationResult) error {

	fileName := fmt.Sprintf("report_%s.md", time.Now().Format("20060102150405"))

	err := os.WriteFile(fileName, []byte(fmt.Sprintf("# Validation Report: %s\n\n%s\n\n----\n# Growth & Discovery\n%s",
		content.EnhancedIdea, content.AuditReport, content.GrowthReport)), 0644)
	if err != nil {
		fmt.Printf("Unable to dump report to file")
		return err
	}

	fmt.Println("\n\n\nSuccessfully dumped report to ", fileName)
	return nil
}

func main() {
	_ = godotenv.Load() // loading .env file

	var rootCmd = &cobra.Command{
		Use:   "validate [args]",
		Short: "A tool to validate your ideas",
		Long: `A tool that helps you validate your business idea or personnel project
	 		to help you understand if it worth building or not and what your competitors are in the sapcve`,
		Run: func(cmd *cobra.Command, args []string) {

			// 1. Launch Bubble Tea Form
			config, err := RunForm()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}

			// 2. Clear output and show selections
			fmt.Printf("🚀 Idea:      %s\n", config.Idea)
			fmt.Printf("🔌 Provider:  %s\n", config.Provider)
			fmt.Printf("🧠 Model:     %s\n", config.Model)
			fmt.Println("--------------------------------------------------")

			result := runValidateCommand(cmd, config)
			err = saveReport(result)
			if err != nil {
				fmt.Println("Dumping report here...\n\n")
				fmt.Println(result)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
