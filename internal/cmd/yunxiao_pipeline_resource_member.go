package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineResourceMemberListOptions struct {
	OrganizationID string
	ResourceType   string
	ResourceID     string
	JSONOutput     bool
}

func newYunxiaoPipelineResourceMemberCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "resource-member",
		Aliases: []string{"resource-members"},
		Short:   "list members who can access a Flow resource",
	}
	command.AddCommand(newYunxiaoPipelineResourceMemberListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoPipelineResourceMemberListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineResourceMemberListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list members with access to a Flow pipeline resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_resource_members", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printPipelineResourceMemberList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ResourceType, "resource-type", "pipeline", "Flow resource type, e.g. pipeline or hostGroup")
	flags.StringVar(&options.ResourceID, "resource-id", "", "Flow resource ID, e.g. pipeline ID")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o pipelineResourceMemberListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"resourceType": strings.TrimSpace(o.ResourceType),
		"resourceId":   strings.TrimSpace(o.ResourceID),
	}
	if params["resourceType"] == "" {
		return nil, fmt.Errorf("resource-type is required")
	}
	if params["resourceId"] == "" {
		return nil, fmt.Errorf("resource-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printPipelineResourceMemberList(out anyWriter, raw string) error {
	rows, ok := pipelineResourceMemberRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tNAME\tROLE")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", row.ID, row.Name, row.Role)
	}
	return writer.Flush()
}

type pipelineResourceMemberRow struct {
	ID   string
	Name string
	Role string
}

func pipelineResourceMemberRowsFromJSON(raw string) []pipelineResourceMemberRow {
	rows, _ := pipelineResourceMemberRowsFromJSONForPrint(raw)
	return rows
}

func pipelineResourceMemberRowsFromJSONForPrint(raw string) ([]pipelineResourceMemberRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]pipelineResourceMemberRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, pipelineResourceMemberRow{
			ID:   firstStringValue(m, "id", "userId", "memberId", "accountId"),
			Name: firstStringValue(m, "name", "displayName", "nickName", "username"),
			Role: firstStringValue(m, "role", "roleName", "roleType", "accessLevel"),
		})
	}
	return rows, true
}
