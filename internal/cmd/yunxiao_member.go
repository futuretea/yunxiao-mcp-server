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
}

func newYunxiaoMemberCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "work with Yunxiao organization members",
	}
	command.AddCommand(newYunxiaoMemberListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoMemberListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options memberListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organization members",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organization_members", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
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
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
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
	_, _ = fmt.Fprintln(writer, "MEMBER_ID\tUSER_ID\tNAME\tEMAIL")
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
