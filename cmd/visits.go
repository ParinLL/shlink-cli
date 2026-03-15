package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var visitCmd = &cobra.Command{
	Use:   "visit",
	Short: "Manage visits",
}

var visitsOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get general visit stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		result, err := shlinkClient.GetVisitsOverview()
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		fmt.Println("Visits Overview:")
		if v := result.Visits.NonOrphanVisits; v != nil {
			fmt.Printf("  Non-orphan: %d (bots: %d)\n", v.Total, v.Bots)
		}
		if v := result.Visits.OrphanVisits; v != nil {
			fmt.Printf("  Orphan:     %d (bots: %d)\n", v.Total, v.Bots)
		}
		return nil
	},
}

var shortURLVisitsCmd = &cobra.Command{
	Use:   "short-url <shortCode>",
	Short: "List visits for a short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		jsonOut, _ := cmd.Flags().GetBool("json")

		result, err := shlinkClient.GetShortURLVisits(args[0], domain, page, perPage)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Referer", "User Agent", "Bot"})
		table.SetAutoWrapText(false)
		for _, v := range result.Visits.Data {
			bot := "No"
			if v.Potentialbot {
				bot = "Yes"
			}
			table.Append([]string{
				v.Date.Format("2006-01-02 15:04"),
				truncate(v.Referer, 30),
				truncate(v.UserAgent, 40),
				bot,
			})
		}
		table.Render()
		p := result.Visits.Pagination
		fmt.Printf("\nPage %d/%d (Total: %d)\n", p.CurrentPage, p.PagesCount, p.TotalItems)
		return nil
	},
}

var tagVisitsCmd = &cobra.Command{
	Use:   "tag <tag>",
	Short: "List visits for a tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		jsonOut, _ := cmd.Flags().GetBool("json")

		result, err := shlinkClient.GetTagVisits(args[0], page, perPage)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Referer", "User Agent", "Bot"})
		table.SetAutoWrapText(false)
		for _, v := range result.Visits.Data {
			bot := "No"
			if v.Potentialbot {
				bot = "Yes"
			}
			table.Append([]string{
				v.Date.Format("2006-01-02 15:04"),
				truncate(v.Referer, 30),
				truncate(v.UserAgent, 40),
				bot,
			})
		}
		table.Render()
		p := result.Visits.Pagination
		fmt.Printf("\nPage %d/%d (Total: %d)\n", p.CurrentPage, p.PagesCount, p.TotalItems)
		return nil
	},
}

var orphanVisitsCmd = &cobra.Command{
	Use:   "orphan",
	Short: "List orphan visits",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		jsonOut, _ := cmd.Flags().GetBool("json")

		result, err := shlinkClient.GetOrphanVisits(page, perPage)
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Referer", "User Agent", "Bot"})
		table.SetAutoWrapText(false)
		for _, v := range result.Visits.Data {
			bot := "No"
			if v.Potentialbot {
				bot = "Yes"
			}
			table.Append([]string{
				v.Date.Format("2006-01-02 15:04"),
				truncate(v.Referer, 30),
				truncate(v.UserAgent, 40),
				bot,
			})
		}
		table.Render()
		p := result.Visits.Pagination
		fmt.Printf("\nPage %d/%d (Total: %d)\n", p.CurrentPage, p.PagesCount, p.TotalItems)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(visitCmd)

	visitsOverviewCmd.Flags().Bool("json", false, "Output as JSON")
	visitCmd.AddCommand(visitsOverviewCmd)

	shortURLVisitsCmd.Flags().String("domain", "", "Domain")
	shortURLVisitsCmd.Flags().Int("page", 1, "Page number")
	shortURLVisitsCmd.Flags().Int("per-page", 20, "Items per page")
	shortURLVisitsCmd.Flags().Bool("json", false, "Output as JSON")
	visitCmd.AddCommand(shortURLVisitsCmd)

	tagVisitsCmd.Flags().Int("page", 1, "Page number")
	tagVisitsCmd.Flags().Int("per-page", 20, "Items per page")
	tagVisitsCmd.Flags().Bool("json", false, "Output as JSON")
	visitCmd.AddCommand(tagVisitsCmd)

	orphanVisitsCmd.Flags().Int("page", 1, "Page number")
	orphanVisitsCmd.Flags().Int("per-page", 20, "Items per page")
	orphanVisitsCmd.Flags().Bool("json", false, "Output as JSON")
	visitCmd.AddCommand(orphanVisitsCmd)
}
