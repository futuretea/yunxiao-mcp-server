package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newYunxiaoTestcaseCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "testcase",
		Aliases: []string{"testcases", "tc"},
		Short:   "work with Projex test cases",
	}
	command.AddCommand(newYunxiaoTestcaseRepoListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoTestcaseViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoTestcaseSearchCommand(streams, cfgFile, v))
	return command
}

// --- repo list ---
type testcaseRepoListOptions struct {
	OrganizationID string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoTestcaseRepoListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options testcaseRepoListOptions
	command := &cobra.Command{
		Use:     "repo-list",
		Aliases: []string{"repos", "repositories"},
		Short:   "list Projex testcase repositories",
		Example: `  # List testcase repositories
  yunxiao testcase repo-list

  # Output as JSON
  yunxiao testcase repo-list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_testcase_repositories", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTestcaseRepoList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o testcaseRepoListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

// --- view ---
type testcaseViewOptions struct {
	OrganizationID string
	TestRepoID     string
	TestcaseID     string
}

func newYunxiaoTestcaseViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options testcaseViewOptions
	command := &cobra.Command{
		Use:     "view <testcase-id>",
		Aliases: []string{"get"},
		Short:   "view a test case as JSON",
		Example: `  # View a test case
  yunxiao testcase view tc-123 --repo-id repo-456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.TestcaseID = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_testcase", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.TestRepoID, "repo-id", "", "testcase repository ID")
	return command
}

func (o testcaseViewOptions) params() (map[string]any, error) {
	params := map[string]any{
		"testRepoId": strings.TrimSpace(o.TestRepoID),
		"testcaseId": strings.TrimSpace(o.TestcaseID),
	}
	if params["testRepoId"] == "" {
		return nil, fmt.Errorf("repo-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

// --- search ---
type testcaseSearchOptions struct {
	OrganizationID string
	TestRepoID     string
	DirectoryID    string
	Subject        string
	OrderBy        string
	Sort           string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoTestcaseSearchCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options testcaseSearchOptions
	command := &cobra.Command{
		Use:   "search",
		Short: "search test cases in a repository",
		Example: `  # Search by subject
  yunxiao testcase search --repo-id repo-456 --subject login

  # Output as JSON
  yunxiao testcase search --repo-id repo-456 --subject login --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "search_testcases", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printTestcaseSearchList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.TestRepoID, "repo-id", "", "testcase repository ID")
	flags.StringVar(&options.DirectoryID, "directory-id", "", "directory ID filter")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field: gmtCreate or name")
	flags.StringVar(&options.Sort, "sort", "", "sort direction: asc or desc")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o testcaseSearchOptions) params() (map[string]any, error) {
	params := map[string]any{
		"testRepoId": strings.TrimSpace(o.TestRepoID),
	}
	if params["testRepoId"] == "" {
		return nil, fmt.Errorf("repo-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "directoryId", o.DirectoryID)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}

// --- print helpers ---

func printTestcaseRepoList(out anyWriter, raw string) error {
	rows, ok := testcaseRepoRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tDESCRIPTION"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Description)
	}
	return writer.Flush()
}

type testcaseRepoRow struct {
	ID          string
	Name        string
	Description string
}

func testcaseRepoRowsFromJSONForPrint(raw string) ([]testcaseRepoRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]testcaseRepoRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, testcaseRepoRow{
			ID:          firstStringValue(m, "id", "testRepoId"),
			Name:        firstStringValue(m, "name", "displayName"),
			Description: firstStringValue(m, "description", "desc"),
		})
	}
	return rows, true
}

func printTestcaseSearchList(out anyWriter, raw string) error {
	rows, ok := testcaseRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tSUBJECT\tPRIORITY"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Subject, row.Priority)
	}
	return writer.Flush()
}

type testcaseRow struct {
	ID       string
	Subject  string
	Priority string
}

func testcaseRowsFromJSONForPrint(raw string) ([]testcaseRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]testcaseRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, testcaseRow{
			ID:       firstStringValue(m, "id", "testcaseId"),
			Subject:  firstStringValue(m, "subject", "title", "name"),
			Priority: firstStringValue(m, "priority", "priorityName"),
		})
	}
	return rows, true
}
