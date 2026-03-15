package cmd

import (
	"fmt"
	"os"

	"shlink-cli/internal/client"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains",
}

var listDomainsCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")
		result, err := shlinkClient.ListDomains()
		if err != nil {
			return err
		}
		if jsonOut {
			return printJSON(result)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Domain", "Default"})
		for _, d := range result.Domains.Data {
			def := "No"
			if d.IsDefault {
				def = "Yes"
			}
			table.Append([]string{d.Domain, def})
		}
		table.Render()
		return nil
	},
}

var setDomainRedirectsCmd = &cobra.Command{
	Use:   "set-redirects <domain>",
	Short: "Set domain 'not found' redirects",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		redirects := &client.DomainRedirects{}
		if cmd.Flags().Changed("base-url-redirect") {
			v, _ := cmd.Flags().GetString("base-url-redirect")
			redirects.BaseUrlRedirect = &v
		}
		if cmd.Flags().Changed("404-redirect") {
			v, _ := cmd.Flags().GetString("404-redirect")
			redirects.Regular404Redirect = &v
		}
		if cmd.Flags().Changed("invalid-redirect") {
			v, _ := cmd.Flags().GetString("invalid-redirect")
			redirects.InvalidShortUrlRedirect = &v
		}

		if err := shlinkClient.SetDomainRedirects(args[0], redirects); err != nil {
			return err
		}
		fmt.Printf("Redirects updated for domain: %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	listDomainsCmd.Flags().Bool("json", false, "Output as JSON")
	domainCmd.AddCommand(listDomainsCmd)

	setDomainRedirectsCmd.Flags().String("base-url-redirect", "", "Base URL redirect")
	setDomainRedirectsCmd.Flags().String("404-redirect", "", "Regular 404 redirect")
	setDomainRedirectsCmd.Flags().String("invalid-redirect", "", "Invalid short URL redirect")
	domainCmd.AddCommand(setDomainRedirectsCmd)
}
