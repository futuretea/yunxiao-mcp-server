package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineJobListOptions struct {
	OrganizationID string
	PipelineID     string
	Category       string
	JSONOutput     bool
}

func newYunxiaoPipelineJobCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "job",
		Aliases: []string{"jobs"},
		Short:   "work with Flow pipeline jobs",
	}
	command.AddCommand(newYunxiaoPipelineJobListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoPipelineJobListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineJobListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Flow pipeline jobs by category",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_pipeline_jobs_by_category", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				_, _ = fmt.Fprintln(streams.Out, result)
				return nil
			}
			return printPipelineJobList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.PipelineID, "pipeline-id", "", "Flow pipeline ID")
	flags.StringVar(&options.Category, "category", "", "task category, e.g. DEPLOY")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o pipelineJobListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"pipelineId": strings.TrimSpace(o.PipelineID),
		"category":   strings.TrimSpace(o.Category),
	}
	if params["pipelineId"] == "" {
		return nil, fmt.Errorf("pipeline-id is required")
	}
	if params["category"] == "" {
		return nil, fmt.Errorf("category is required, e.g. DEPLOY")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
}

func printPipelineJobList(out anyWriter, raw string) error {
	rows, ok := pipelineJobRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "IDENTIFIER\tNAME\tCATEGORY\tSTATUS")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.Identifier, row.Name, row.Category, row.Status)
	}
	return writer.Flush()
}

type pipelineJobRow struct {
	Identifier string
	Name       string
	Category   string
	Status     string
}

func pipelineJobRowsFromJSON(raw string) []pipelineJobRow {
	rows, _ := pipelineJobRowsFromJSONForPrint(raw)
	return rows
}

func pipelineJobRowsFromJSONForPrint(raw string) ([]pipelineJobRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]pipelineJobRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, pipelineJobRow{
			Identifier: firstStringValue(m, "identifier", "id", "jobId"),
			Name:       firstStringValue(m, "name", "displayName", "taskName"),
			Category:   firstStringValue(m, "category", "taskCategory"),
			Status:     firstStringValue(m, "status", "taskStatus", "state"),
		})
	}
	return rows, true
}
