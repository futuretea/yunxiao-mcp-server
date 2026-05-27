package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type memberListOptions struct {
	OrganizationID string
	Page           int
	PerPage        int
	JSONOutput     bool
	OutputFormat string
}

func newYunxiaoMemberCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "work with Yunxiao organization members",
	}
	command.AddCommand(newYunxiaoMemberListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoMemberSearchCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoMemberListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options memberListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organization members",
		Example: `  # List members (uses default organization)
  yunxiao member list

  # List members from a specific organization
  yunxiao member list --organization-id org-abc

  # Search by name or email
  yunxiao member list --organization-id org-abc --query alice`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organization_members", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printMemberList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o memberListOptions) params() map[string]any {
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

func printMemberList(out anyWriter, raw string) error {
	rows := memberRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("MEMBER_ID\tUSER_ID\tNAME\tEMAIL"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.MemberID, row.UserID, row.Name, row.Email)
	}
	return writer.Flush()
}

type memberRow struct {
	MemberID string
	UserID   string
	Name     string
	Email    string
}

func memberRowsFromJSON(raw string) []memberRow {
	items := rowsFromJSON(raw)
	rows := make([]memberRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, memberRow{
			MemberID: firstStringValue(m, "id", "identifier", "memberId"),
			UserID:   firstStringValue(m, "userId", "accountId"),
			Name:     firstStringValue(m, "displayName", "name", "username", "nickName", "realName"),
			Email:    firstStringValue(m, "email", "mail", "emailAddress"),
		})
	}
	return rows
}

type memberSearchOptions struct {
	OrganizationID  string
	Query           string
	DeptIDs         string
	RoleIDs         string
	Statuses        string
	IncludeChildren bool
	Page            int
	PerPage         int
	JSONOutput      bool
	OutputFormat string
}

func newYunxiaoMemberSearchCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options memberSearchOptions
	command := &cobra.Command{
		Use:   "search",
		Short: "search Yunxiao organization members",
		Example: `  # Search by name
  yunxiao member search --query alice

  # Search with department filter
  yunxiao member search --query alice --dept-ids "dept-1,dept-2"

  # Include child departments
  yunxiao member search --query alice --include-children

  # Output as JSON
  yunxiao member search --query alice --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "search_organization_members", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput || options.OutputFormat == "json" {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printMemberList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.Query, "query", "", "search keyword for username, display name, or email")
	flags.StringVar(&options.DeptIDs, "dept-ids", "", "comma-separated department IDs")
	flags.StringVar(&options.RoleIDs, "role-ids", "", "comma-separated role IDs")
	flags.StringVar(&options.Statuses, "statuses", "", "comma-separated member statuses, e.g. enabled,disabled")
	flags.BoolVar(&options.IncludeChildren, "include-children", false, "include members from child departments")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	flags.StringVar(&options.OutputFormat, "output", "", "output format: table, json, or csv")
	return command
}

func (o memberSearchOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "query", o.Query)
	setCLIStringParam(params, "deptIds", o.DeptIDs)
	setCLIStringParam(params, "roleIds", o.RoleIDs)
	setCLIStringParam(params, "statuses", o.Statuses)
	if o.IncludeChildren {
		params["includeChildren"] = true
	}
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}
