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
	command.AddCommand(newYunxiaoProjectMemberCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectRoleCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoProjectViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoProjectListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options projectListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Projex projects",
		Example: `  # List all projects
  yunxiao project list

  # Search by name keyword
  yunxiao project list --name demo

  # Output as JSON
  yunxiao project list --json`,
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
				printCLIJSON(streams.Out, result)
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
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
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
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tNAME\tSTATUS\tCREATOR"))
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

type projectViewOptions struct {
	OrganizationID    string
	ProjectID         string
	IncludeMembers    bool
	IncludeSprints    bool
	IncludeMilestones bool
	IncludeVersions   bool
	IncludeLabels     bool
	ActiveOnly        bool
	Status            string
	Page              int
	PerPage           int
}

func newYunxiaoProjectViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := projectViewOptions{
		IncludeMembers:    true,
		IncludeSprints:    true,
		IncludeMilestones: true,
		IncludeVersions:   true,
		IncludeLabels:     true,
		ActiveOnly:        true,
	}
	command := &cobra.Command{
		Use:     "view <project-id>",
		Aliases: []string{"overview", "info"},
		Short:   "view a Projex project overview as JSON",
		Example: `  # View project overview
  yunxiao project view 123

  # View without sprint details
  yunxiao project view 123 --include-sprints=false`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.ProjectID = args[0]
			result, err := callYunxiaoTool(cmd, cfg, "get_project_overview", options.params())
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeMembers, "include-members", true, "include project members")
	flags.BoolVar(&options.IncludeSprints, "include-sprints", true, "include sprints list")
	flags.BoolVar(&options.IncludeMilestones, "include-milestones", true, "include milestones list")
	flags.BoolVar(&options.IncludeVersions, "include-versions", true, "include versions list")
	flags.BoolVar(&options.IncludeLabels, "include-labels", true, "include labels list")
	flags.BoolVar(&options.ActiveOnly, "active-only", true, "show only active sprints, milestones, and versions")
	flags.StringVar(&options.Status, "status", "", "comma-separated statuses for active-only filter; defaults to TODO,DOING")
	flags.IntVar(&options.Page, "page", 0, "page number for paged sections")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size for paged sections")
	return command
}

func (o projectViewOptions) params() map[string]any {
	params := map[string]any{
		"projectId":         o.ProjectID,
		"includeMembers":    o.IncludeMembers,
		"includeSprints":    o.IncludeSprints,
		"includeMilestones": o.IncludeMilestones,
		"includeVersions":   o.IncludeVersions,
		"includeLabels":     o.IncludeLabels,
		"activeOnly":        o.ActiveOnly,
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "status", o.Status)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}
