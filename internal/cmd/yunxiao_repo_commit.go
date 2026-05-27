package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoCommitListOptions struct {
	OrganizationID   string
	RepositoryID     string
	Ref              string
	Since            string
	Until            string
	Page             int
	PerPage          int
	Path             string
	Search           string
	ShowSignature    bool
	ShowSignatureSet bool
	CommitterIDs     string
	JSONOutput       bool
	OutputFormat string
}

func newYunxiaoRepoCommitCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "commit",
		Aliases: []string{"commits"},
		Short:   "work with CodeUp repository commits",
	}
	command.AddCommand(newYunxiaoRepoCommitListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoCommitViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoCommitListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoCommitListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list CodeUp repository commits",
		Example: `  # List commits
  yunxiao repo commit list --repository-id group/repo

  # List commits on a specific branch
  yunxiao repo commit list --repository-id group/repo --ref-name main

  # Output as JSON
  yunxiao repo commit list --repository-id group/repo --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ShowSignatureSet = cmd.Flags().Changed("show-signature")
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_commits", params)
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRepoCommitList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "CodeUp repository ID or full path")
	flags.StringVar(&options.Ref, "ref", "", "branch, tag, or commit SHA")
	flags.StringVar(&options.Since, "since", "", "start time in ISO 8601 format")
	flags.StringVar(&options.Until, "until", "", "end time in ISO 8601 format")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.StringVar(&options.Path, "path", "", "filter commits touching this path")
	flags.StringVar(&options.Search, "search", "", "commit search keyword")
	flags.BoolVar(&options.ShowSignature, "show-signature", false, "include commit signatures; use --show-signature=false to exclude explicitly")
	flags.StringVar(&options.CommitterIDs, "committer-ids", "", "comma-separated committer user IDs")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o repoCommitListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"refName":      strings.TrimSpace(o.Ref),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["refName"] == "" {
		return nil, fmt.Errorf("ref is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "since", o.Since)
	setCLIStringParam(params, "until", o.Until)
	setCLIStringParam(params, "path", o.Path)
	setCLIStringParam(params, "search", o.Search)
	setCLIStringParam(params, "committerIds", o.CommitterIDs)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	if o.ShowSignatureSet {
		params["showSignature"] = o.ShowSignature
	}
	return params, nil
}

func printRepoCommitList(out anyWriter, raw string) error {
	rows, ok := repoCommitRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("SHA\tSHORT_ID\tTITLE\tAUTHOR\tDATE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", row.SHA, row.ShortID, row.Title, row.Author, row.Date)
	}
	return writer.Flush()
}

type repoCommitRow struct {
	SHA     string
	ShortID string
	Title   string
	Author  string
	Date    string
}


func repoCommitRowsFromJSONForPrint(raw string) ([]repoCommitRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]repoCommitRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, repoCommitRow{
			SHA:     firstStringValue(m, "id", "sha", "commitId"),
			ShortID: firstStringValue(m, "shortId", "short_id"),
			Title:   firstStringValue(m, "title", "message", "subject"),
			Author:  firstStringValue(m, "authorName", "author_name", "author", "committerName", "committer"),
			Date:    firstStringValue(m, "authoredDate", "authored_date", "committedDate", "committed_date", "createdAt"),
		})
	}
	return rows, true
}
