package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectListOptions struct {
	OrganizationID string
	Name           string
	Status         string
	Creator        string
	OrderBy        string
	Sort           string
	Page           int
	PerPage        int
	JSONOutput     bool
}

func newYunxiaoProjectCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects"},
		Short:   "work with Projex projects",
	}
	command.AddCommand(newYunxiaoProjectListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoProjectListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "search_projects", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				_, _ = fmt.Fprintln(streams.Out, result)
				return nil
			}
			return printProjectList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Name, "name", "", "project name keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated project status IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.StringVar(&options.OrderBy, "order-by", "", "sort field")
	flags.StringVar(&options.Sort, "sort", "", "sort direction, e.g. asc or desc")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o projectListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "name", o.Name)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "creator", o.Creator)
	setCLIStringParam(params, "orderBy", o.OrderBy)
	setCLIStringParam(params, "sort", o.Sort)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printProjectList(out anyWriter, raw string) error {
	rows := projectRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tNAME\tSTATUS\tCREATOR")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Name, row.Status, row.Creator)
	}
	return writer.Flush()
}

type projectRow struct {
	ID      string
	Name    string
	Status  string
	Creator string
}

func projectRowsFromJSON(raw string) []projectRow {
	items := rowsFromJSON(raw)
	rows := make([]projectRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, projectRow{
			ID:      firstStringValue(m, "id", "identifier", "projectId", "spaceId"),
			Name:    firstStringValue(m, "name", "title"),
			Status:  firstStringValue(m, "status", "statusName", "statusIdentifier"),
			Creator: firstStringValue(m, "creator", "owner", "creatorName", "ownerName"),
		})
	}
	return rows
}
