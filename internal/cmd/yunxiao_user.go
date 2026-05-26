package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type userListOptions struct {
	Filter     string
	Status     string
	DeptID     string
	Page       int
	PerPage    int
	JSONOutput bool
}

type userViewOptions struct {
	IDOrUsername  string
	OrganizationID string
}

func newYunxiaoUserCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		Short:   "work with Yunxiao users",
	}
	command.AddCommand(newYunxiaoUserWhoamiCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoUserListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoUserGetCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoUserWhoamiCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:   "whoami",
		Short: "show the current authenticated user as JSON",
		Example: `  # Show current user
  yunxiao user whoami`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_current_user", map[string]any{})
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	return command
}

func newYunxiaoUserListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options userListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao users",
		Example: `  # List users (uses default organization)
  yunxiao user list

  # Search by name or email
  yunxiao user list --filter alice

  # Filter by department
  yunxiao user list --dept-id dept-123

  # Output as JSON
  yunxiao user list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_users", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printUserList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.Filter, "filter", "", "fuzzy filter for username, login, email, or phone")
	flags.StringVar(&options.Status, "status", "", "user status, e.g. enabled or deleted")
	flags.StringVar(&options.DeptID, "dept-id", "", "department ID filter")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.IntVar(&options.PerPage, "limit", 0, "max results (alias for --per-page)")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func newYunxiaoUserGetCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options userViewOptions
	command := &cobra.Command{
		Use:     "get <id-or-username>",
		Aliases: []string{"view"},
		Short:   "get a Yunxiao user as JSON",
		Example: `  # Get user by ID
  yunxiao user get user-12345

  # Get user by username
  yunxiao user get alice`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.IDOrUsername = args[0]
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_user", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	return command
}

func (o userListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "filter", o.Filter)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "deptId", o.DeptID)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func (o userViewOptions) params() (map[string]any, error) {
	params := map[string]any{
		"idOrUsername": strings.TrimSpace(o.IDOrUsername),
	}
	if params["idOrUsername"] == "" {
		return nil, fmt.Errorf("id-or-username is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printUserList(out anyWriter, raw string) error {
	rows, ok := userRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("ID\tUSERNAME\tEMAIL\tSTATUS"))
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Username, row.Email, row.Status)
	}
	return writer.Flush()
}

type userRow struct {
	ID       string
	Username string
	Email    string
	Status   string
}

func userRowsFromJSONForPrint(raw string) ([]userRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]userRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, userRow{
			ID:       firstStringValue(m, "id", "userId", "accountId"),
			Username: firstStringValue(m, "username", "login", "name"),
			Email:    firstStringValue(m, "email"),
			Status:   firstStringValue(m, "status"),
		})
	}
	return rows, true
}
