package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check Shlink instance health",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		result, err := shlinkClient.Health()
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		fmt.Printf("Status:  %s\n", result.Status)
		fmt.Printf("Version: %s\n", result.Version)
		return nil
	},
}

func init() {
	healthCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(healthCmd)
}
