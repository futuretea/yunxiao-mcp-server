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
}

func newYunxiaoSprintCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "sprint",
		Aliases: []string{"sprints"},
		Short:   "work with Projex sprints",
	}
	command.AddCommand(newYunxiaoSprintListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoSprintViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoSprintListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options sprintListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex sprints in a project",
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
			if options.JSONOutput {
				_, _ = fmt.Fprintln(streams.Out, result)
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
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
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
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tNAME\tSTATUS")
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
