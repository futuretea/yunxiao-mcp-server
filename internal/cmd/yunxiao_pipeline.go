package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineListOptions struct {
	OrganizationID   string
	PipelineName     string
	StatusList       string
	CreateStartTime  int64
	CreateEndTime    int64
	ExecuteStartTime int64
	ExecuteEndTime   int64
	Page             int
	PerPage          int
	JSONOutput       bool
}

func newYunxiaoPipelineCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipelines"},
		Short:   "work with Flow pipelines",
	}
	command.AddCommand(newYunxiaoPipelineListCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoPipelineListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Flow pipelines",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_pipelines", options.params())
			if err != nil {
				return err
			}
			if options.JSONOutput {
				_, _ = fmt.Fprintln(streams.Out, result)
				return nil
			}
			return printPipelineList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.PipelineName, "name", "", "pipeline name keyword")
	flags.StringVar(&options.StatusList, "status", "", "comma-separated pipeline statuses, e.g. RUNNING,SUCCESS,FAIL")
	flags.Int64Var(&options.CreateStartTime, "create-start-time", 0, "pipeline creation start time as Unix milliseconds")
	flags.Int64Var(&options.CreateEndTime, "create-end-time", 0, "pipeline creation end time as Unix milliseconds")
	flags.Int64Var(&options.ExecuteStartTime, "execute-start-time", 0, "pipeline execution start time as Unix milliseconds")
	flags.Int64Var(&options.ExecuteEndTime, "execute-end-time", 0, "pipeline execution end time as Unix milliseconds")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o pipelineListOptions) params() map[string]any {
	params := map[string]any{}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "pipelineName", o.PipelineName)
	setCLIStringParam(params, "statusList", o.StatusList)
	if o.CreateStartTime > 0 {
		params["createStartTime"] = o.CreateStartTime
	}
	if o.CreateEndTime > 0 {
		params["createEndTime"] = o.CreateEndTime
	}
	if o.ExecuteStartTime > 0 {
		params["executeStartTime"] = o.ExecuteStartTime
	}
	if o.ExecuteEndTime > 0 {
		params["executeEndTime"] = o.ExecuteEndTime
	}
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params
}

func printPipelineList(out anyWriter, raw string) error {
	rows, ok := pipelineRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tNAME\tSTATUS\tLAST_RUN")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", row.ID, row.Name, row.Status, row.LastRun)
	}
	return writer.Flush()
}

type pipelineRow struct {
	ID      string
	Name    string
	Status  string
	LastRun string
}

func pipelineRowsFromJSON(raw string) []pipelineRow {
	rows, _ := pipelineRowsFromJSONForPrint(raw)
	return rows
}

func pipelineRowsFromJSONForPrint(raw string) ([]pipelineRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]pipelineRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, pipelineRow{
			ID:      firstStringValue(m, "pipelineId", "id", "pipelineID"),
			Name:    firstStringValue(m, "pipelineName", "name", "displayName"),
			Status:  firstStringValue(m, "status", "latestStatus", "lastRunStatus"),
			LastRun: pipelineLastRunValue(m),
		})
	}
	return rows, true
}

func pipelineLastRunValue(m map[string]any) string {
	for _, key := range []string{"latestRun", "lastRun", "pipelineRun", "latestPipelineRun"} {
		if nested, ok := m[key].(map[string]any); ok {
			if value := firstStringValue(nested, "pipelineRunId", "runId", "id", "status", "runStatus"); value != "" {
				return value
			}
		}
	}
	return firstStringValue(m, "latestRunId", "lastRunId", "pipelineRunId", "lastPipelineRunId")
}
