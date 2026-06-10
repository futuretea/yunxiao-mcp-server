package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type sprintListOptions struct {
	OrganizationID string
	ProjectID      string
	Status         string
	Name           string
	Page           int
	PerPage        int
	JSONOutput     bool
	OutputFormat   string
}

func newYunxiaoSprintCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "sprint",
		Aliases: []string{"sprints"},
		Short:   "work with Projex sprints",
	}
	command.AddCommand(newYunxiaoSprintListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoSprintViewCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoSprintVelocityCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoSprintListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options sprintListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex sprints in a project",
		Example: `  # List active sprints
  yunxiao sprint list --project-id 123

  # List all sprints including archived
  yunxiao sprint list --project-id 123 --status "TODO,DOING,ARCHIVED"

  # Output as JSON
  yunxiao sprint list --project-id 123 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_sprints", params)
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printSprintList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.StringVar(&options.Status, "status", "", "comma-separated sprint statuses, e.g. TODO,DOING,ARCHIVED")
	flags.StringVar(&options.Name, "name", "", "sprint name keyword")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o sprintListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"projectId": strings.TrimSpace(o.ProjectID),
	}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "name", o.Name)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}

func printSprintList(out anyWriter, raw string) error {
	rows := sprintRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tSTATUS"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Status)
	}
	return writer.Flush()
}

type sprintRow struct {
	ID     string
	Name   string
	Status string
}

func sprintRowsFromJSON(raw string) []sprintRow {
	items := rowsFromJSON(raw)
	rows := make([]sprintRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, sprintRow{
			ID:     firstStringValue(m, "id", "identifier", "sprintId"),
			Name:   firstStringValue(m, "name", "title"),
			Status: firstStringValue(m, "status", "statusName", "statusIdentifier"),
		})
	}
	return rows
}

type sprintVelocityOptions struct {
	OrganizationID string
	ProjectID      string
	Categories     string
	SprintCount    int
	SprintStatus   string
}

func newYunxiaoSprintVelocityCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options sprintVelocityOptions
	command := &cobra.Command{
		Use:     "velocity <project-id>",
		Aliases: []string{"metrics"},
		Short:   "show sprint velocity metrics as JSON",
		Example: `  # View velocity for last 5 sprints
  yunxiao sprint velocity 123

  # View velocity for last 10 sprints
  yunxiao sprint velocity 123 --sprint-count 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_sprint_velocity", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated categories; defaults to Task,Bug")
	flags.IntVar(&options.SprintCount, "sprint-count", 0, "number of recent sprints to analyze; defaults to 5, max 20")
	flags.StringVar(&options.SprintStatus, "sprint-status", "", "comma-separated sprint statuses; defaults to ARCHIVED,DONE")
	return command
}

func (o sprintVelocityOptions) params() map[string]any {
	params := map[string]any{"projectId": o.ProjectID}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "sprintStatus", o.SprintStatus)
	if o.SprintCount > 0 {
		params["sprintCount"] = o.SprintCount
	}
	return params
}
