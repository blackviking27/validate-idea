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

	result, err := agents.RunValidator(ctx, aiProvider, userIdea)
	if err != nil {
		log.Fatalf("Validation failed : %v", err)
	}

	return result
}

func saveReport(content string) error {
	fileName := fmt.Sprintf("report_%s.md", time.Now().Format("20060102150405"))

	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Unable to dump report to file")
		return err
	}

	fmt.Println("Successfully dumped report to ", fileName)
	return nil
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
			err := saveReport(result)
			if err != nil {
				fmt.Println("Dumping report here...\n\n")
				fmt.Println(result)
			}
		},
	}

	//adding flags
	rootCmd.PersistentFlags().StringVarP(&providerArg, "provider", "p", "", "AI model provider")
	rootCmd.PersistentFlags().StringVarP(&aiModel, "model", "m", "", "AI model to use")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
