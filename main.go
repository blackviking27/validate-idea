package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/blackviking27/validate-idea-cli/agents"
	"github.com/blackviking27/validate-idea-cli/providers"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	providerArg string
)

func getAiProvider(ctx context.Context) *providers.AIProvider {
	var currProivder providers.AIProvider

	switch providerArg {
	case "gemini":
		currProivder, _ = providers.NewGeminiProvider(ctx)
	case "ollama":
		currProivder, _ = providers.NewOllamaProvider(ctx)
	}

	if currProivder == nil {
		log.Fatal("Not a valid AI provider provided")
	}

	return &currProivder
}

func runValidateCommand(cmd *cobra.Command, args []string) string {
	ctx := context.Background()
	userIdea := args[0]

	// Initiating ai provider
	aiProvider := getAiProvider(ctx)

	fmt.Printf("🔍 Analyzing idea with %s...\n\n", (*aiProvider).Name())

	result, err := agents.RunValidator(ctx, aiProvider, userIdea)
	if err != nil {
		log.Fatalf("Validation failed : %v", err)
	}

	return result

}

func main() {
	_ = godotenv.Load() // loading .env file

	var rootCmd = &cobra.Command{
		Use:   "validate [args]",
		Short: "A tool to validate your ideas",
		Long: `A tool that helps you validate your business idea or personnel project
	 		to help you understand if it worth building or not and what your competitors are in the sapcve`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			result := runValidateCommand(cmd, args)
			fmt.Println(result)
		},
	}

	//adding flags
	rootCmd.PersistentFlags().StringVarP(&providerArg, "provider", "p", "", "AI model provider")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
