package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type organizationListOptions struct {
	Page       int
	PerPage    int
	JSONOutput bool
}

func newYunxiaoOrganizationCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "organization",
		Aliases: []string{"organizations"},
		Short:   "work with Yunxiao organizations",
	}
	command.AddCommand(newYunxiaoOrganizationListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoOrganizationListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options organizationListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Yunxiao organizations visible to the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_organizations", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				printCLIJSON(streams.Out, result)
				return nil
			}
			return printOrganizationList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o organizationListOptions) params() map[string]any {
	params := map[string]any{}
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printOrganizationList(out anyWriter, raw string) error {
	rows := organizationRowsFromJSON(raw)
	if len(rows) == 0 {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tNAME")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", row.ID, row.Name)
	}
	return writer.Flush()
}

type organizationRow struct {
	ID   string
	Name string
}

func organizationRowsFromJSON(raw string) []organizationRow {
	items := rowsFromJSON(raw)
	rows := make([]organizationRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, organizationRow{
			ID:   firstStringValue(m, "id", "identifier", "organizationId"),
			Name: firstStringValue(m, "name", "displayName", "title"),
		})
	}
	return rows
}
