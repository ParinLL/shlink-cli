package cmd

import (
	"fmt"
	"os"

	"shlink-cli/internal/client"

	"github.com/spf13/cobra"
)

var (
	debug      bool
	baseURL    string
	apiKey     string
	shlinkClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "shlink-cli",
	Short: "CLI client for Shlink URL shortener",
	Long:  "A command-line interface for managing Shlink, the self-hosted URL shortener.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if baseURL == "" {
			baseURL = os.Getenv("SHLINK_BASE_URL")
		}
		if apiKey == "" {
			apiKey = os.Getenv("SHLINK_API_KEY")
		}

		// Skip validation for health and help commands
		if cmd.Name() == "health" || cmd.Name() == "help" || cmd.Name() == "version" {
			shlinkClient = client.New(baseURL, apiKey, debug)
			return
		}

		if baseURL == "" {
			fmt.Fprintln(os.Stderr, "Error: SHLINK_BASE_URL environment variable or --base-url flag is required")
			os.Exit(1)
		}
		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "Error: SHLINK_API_KEY environment variable or --api-key flag is required")
			os.Exit(1)
		}

		shlinkClient = client.New(baseURL, apiKey, debug)

		if debug {
			fmt.Printf("[DEBUG] Base URL: %s\n", baseURL)
			fmt.Println("[DEBUG] API Key: [REDACTED]")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output (API key will be redacted)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Shlink base URL (or set SHLINK_BASE_URL)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Shlink API key (or set SHLINK_API_KEY)")
}
