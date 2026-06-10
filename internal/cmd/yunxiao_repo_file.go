package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoFileListOptions struct {
	OrganizationID string
	RepositoryID   string
	Path           string
	Ref            string
	TreeType       string
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoRepoFileCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "file",
		Aliases: []string{"files"},
		Short:   "work with CodeUp repository files",
	}
	command.AddCommand(newYunxiaoRepoFileListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoRepoFileViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoRepoFileListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoFileListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list CodeUp repository files",
		Example: `  # List files in root
  yunxiao repo file list --repository-id group/repo

  # List files in a subdirectory
  yunxiao repo file list --repository-id group/repo --path src/

  # List files on a specific branch
  yunxiao repo file list --repository-id group/repo --ref-name develop

  # Output as JSON
  yunxiao repo file list --repository-id group/repo --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_files", params)
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printRepoFileList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "CodeUp repository ID or full path")
	flags.StringVar(&options.Path, "path", "", "directory path to query")
	flags.StringVar(&options.Ref, "ref", "", "branch, tag, or commit SHA")
	flags.StringVar(&options.TreeType, "type", "", "tree mode: DIRECT, RECURSIVE, or FLATTEN")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o repoFileListOptions) params() (map[string]any, error) {
	params := map[string]any{"repositoryId": strings.TrimSpace(o.RepositoryID)}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "path", o.Path)
	setCLIStringParam(params, "ref", o.Ref)
	setCLIStringParam(params, "type", o.TreeType)
	return params, nil
}

func printRepoFileList(out anyWriter, raw string) error {
	rows, ok := repoFileRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("PATH\tTYPE\tSIZE\tMODE"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.Path, row.Type, row.Size, row.Mode)
	}
	return writer.Flush()
}

type repoFileRow struct {
	Path string
	Type string
	Size string
	Mode string
}

func repoFileRowsFromJSONForPrint(raw string) ([]repoFileRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]repoFileRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, repoFileRow{
			Path: firstStringValue(m, "path", "filePath", "name"),
			Type: firstStringValue(m, "type", "kind", "fileType"),
			Size: firstStringValue(m, "size", "fileSize"),
			Mode: firstStringValue(m, "mode"),
		})
	}
	return rows, true
}
