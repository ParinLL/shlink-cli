package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ParinLL/shlink-cli/internal/client"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var shortURLCmd = &cobra.Command{
	Use:     "short-url",
	Aliases: []string{"su"},
	Short:   "Manage short URLs",
}

var listShortURLsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all short URLs",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		search, _ := cmd.Flags().GetString("search")
		orderBy, _ := cmd.Flags().GetString("order-by")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		jsonOut, _ := cmd.Flags().GetBool("json")

		result, err := shlinkClient.ListShortURLs(page, perPage, search, orderBy, tags)
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(result)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Short Code", "Short URL", "Long URL", "Tags", "Visits"})
		table.SetAutoWrapText(false)

		for _, su := range result.ShortUrls.Data {
			visits := "0"
			if su.VisitsSummary != nil {
				visits = fmt.Sprintf("%d", su.VisitsSummary.Total)
			}
			table.Append([]string{
				su.ShortCode,
				su.ShortURL,
				truncate(su.LongURL, 50),
				strings.Join(su.Tags, ", "),
				visits,
			})
		}
		table.Render()

		p := result.ShortUrls.Pagination
		fmt.Printf("\nPage %d/%d (Total: %d)\n", p.CurrentPage, p.PagesCount, p.TotalItems)
		return nil
	},
}

var createShortURLCmd = &cobra.Command{
	Use:   "create <longUrl>",
	Short: "Create a new short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug, _ := cmd.Flags().GetString("slug")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		title, _ := cmd.Flags().GetString("title")
		domain, _ := cmd.Flags().GetString("domain")
		findIfExists, _ := cmd.Flags().GetBool("find-if-exists")
		maxVisits, _ := cmd.Flags().GetInt("max-visits")
		jsonOut, _ := cmd.Flags().GetBool("json")

		req := &client.CreateShortURLRequest{
			LongURL:      args[0],
			CustomSlug:   slug,
			Tags:         tags,
			Title:        title,
			Domain:       domain,
			FindIfExists: findIfExists,
			MaxVisits:    maxVisits,
		}

		result, err := shlinkClient.CreateShortURL(req)
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(result)
		}

		fmt.Printf("Short URL created:\n")
		fmt.Printf("  Short Code: %s\n", result.ShortCode)
		fmt.Printf("  Short URL:  %s\n", result.ShortURL)
		fmt.Printf("  Long URL:   %s\n", result.LongURL)
		return nil
	},
}

var getShortURLCmd = &cobra.Command{
	Use:   "get <shortCode>",
	Short: "Get details of a short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		jsonOut, _ := cmd.Flags().GetBool("json")

		result, err := shlinkClient.GetShortURL(args[0], domain)
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(result)
		}

		fmt.Printf("Short Code:  %s\n", result.ShortCode)
		fmt.Printf("Short URL:   %s\n", result.ShortURL)
		fmt.Printf("Long URL:    %s\n", result.LongURL)
		fmt.Printf("Created:     %s\n", result.DateCreated.Format("2006-01-02 15:04:05"))
		fmt.Printf("Tags:        %s\n", strings.Join(result.Tags, ", "))
		if result.Title != nil {
			fmt.Printf("Title:       %s\n", *result.Title)
		}
		if result.VisitsSummary != nil {
			fmt.Printf("Visits:      %d (bots: %d)\n", result.VisitsSummary.Total, result.VisitsSummary.Bots)
		}
		return nil
	},
}

var editShortURLCmd = &cobra.Command{
	Use:   "edit <shortCode>",
	Short: "Edit an existing short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		jsonOut, _ := cmd.Flags().GetBool("json")

		req := &client.EditShortURLRequest{}
		if cmd.Flags().Changed("long-url") {
			v, _ := cmd.Flags().GetString("long-url")
			req.LongURL = &v
		}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			req.Title = &v
		}
		if cmd.Flags().Changed("tags") {
			v, _ := cmd.Flags().GetStringSlice("tags")
			req.Tags = v
		}
		if cmd.Flags().Changed("max-visits") {
			v, _ := cmd.Flags().GetInt("max-visits")
			req.MaxVisits = &v
		}

		result, err := shlinkClient.EditShortURL(args[0], domain, req)
		if err != nil {
			return err
		}

		if jsonOut {
			return printJSON(result)
		}

		fmt.Printf("Short URL updated: %s -> %s\n", result.ShortURL, result.LongURL)
		return nil
	},
}

var deleteShortURLCmd = &cobra.Command{
	Use:   "delete <shortCode>",
	Short: "Delete a short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")

		if err := shlinkClient.DeleteShortURL(args[0], domain); err != nil {
			return err
		}

		fmt.Printf("Short URL '%s' deleted.\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(shortURLCmd)

	// list
	listShortURLsCmd.Flags().Int("page", 1, "Page number")
	listShortURLsCmd.Flags().Int("per-page", 20, "Items per page")
	listShortURLsCmd.Flags().String("search", "", "Search term")
	listShortURLsCmd.Flags().String("order-by", "", "Order by field (e.g. dateCreated-DESC)")
	listShortURLsCmd.Flags().StringSlice("tags", nil, "Filter by tags")
	listShortURLsCmd.Flags().Bool("json", false, "Output as JSON")
	shortURLCmd.AddCommand(listShortURLsCmd)

	// create
	createShortURLCmd.Flags().String("slug", "", "Custom slug")
	createShortURLCmd.Flags().StringSlice("tags", nil, "Tags")
	createShortURLCmd.Flags().String("title", "", "Title")
	createShortURLCmd.Flags().String("domain", "", "Domain")
	createShortURLCmd.Flags().Bool("find-if-exists", false, "Return existing if same long URL exists")
	createShortURLCmd.Flags().Int("max-visits", 0, "Max visits (0 = unlimited)")
	createShortURLCmd.Flags().Bool("json", false, "Output as JSON")
	shortURLCmd.AddCommand(createShortURLCmd)

	// get
	getShortURLCmd.Flags().String("domain", "", "Domain")
	getShortURLCmd.Flags().Bool("json", false, "Output as JSON")
	shortURLCmd.AddCommand(getShortURLCmd)

	// edit
	editShortURLCmd.Flags().String("domain", "", "Domain")
	editShortURLCmd.Flags().String("long-url", "", "New long URL")
	editShortURLCmd.Flags().String("title", "", "New title")
	editShortURLCmd.Flags().StringSlice("tags", nil, "New tags")
	editShortURLCmd.Flags().Int("max-visits", 0, "Max visits")
	editShortURLCmd.Flags().Bool("json", false, "Output as JSON")
	shortURLCmd.AddCommand(editShortURLCmd)

	// delete
	deleteShortURLCmd.Flags().String("domain", "", "Domain")
	shortURLCmd.AddCommand(deleteShortURLCmd)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
