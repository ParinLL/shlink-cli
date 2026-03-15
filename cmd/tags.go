package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var listTagsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, _ := cmd.Flags().GetBool("stats")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if stats {
			result, err := shlinkClient.TagsWithStats()
			if err != nil {
				return err
			}
			if jsonOut {
				return printJSON(result)
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Tag", "Short URLs", "Visits"})
			for _, t := range result.Tags.Data {
				visits := "0"
				if t.VisitsSummary != nil {
					visits = fmt.Sprintf("%d", t.VisitsSummary.Total)
				}
				table.Append([]string{t.Tag, fmt.Sprintf("%d", t.ShortUrlsCount), visits})
			}
			table.Render()
			return nil
		}

		result, err := shlinkClient.ListTags()
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		for _, tag := range result.Tags.Data {
			fmt.Println(tag)
		}
		return nil
	},
}

var renameTagCmd = &cobra.Command{
	Use:   "rename <oldName> <newName>",
	Short: "Rename a tag",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := shlinkClient.RenameTag(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Tag renamed: %s -> %s\n", args[0], args[1])
		return nil
	},
}

var deleteTagsCmd = &cobra.Command{
	Use:   "delete <tag1> [tag2...]",
	Short: "Delete one or more tags",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := shlinkClient.DeleteTags(args); err != nil {
			return err
		}
		fmt.Printf("Deleted %d tag(s).\n", len(args))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	listTagsCmd.Flags().Bool("stats", false, "Include visit stats")
	listTagsCmd.Flags().Bool("json", false, "Output as JSON")
	tagCmd.AddCommand(listTagsCmd)

	tagCmd.AddCommand(renameTagCmd)
	tagCmd.AddCommand(deleteTagsCmd)
}
