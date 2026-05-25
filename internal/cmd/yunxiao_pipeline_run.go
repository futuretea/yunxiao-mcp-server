package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineRunListOptions struct {
	OrganizationID string
	PipelineID     string
	Page           int
	PerPage        int
	StartTime      int64
	EndTime        int64
	Status         string
	TriggerMode    int
	JSONOutput     bool
}

func newYunxiaoPipelineRunCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "run",
		Aliases: []string{"runs"},
		Short:   "work with Flow pipeline runs",
	}
	command.AddCommand(newYunxiaoPipelineRunListCommand(streams, cfgFile, v))
	command.AddCommand(newYunxiaoPipelineRunViewCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoPipelineRunListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineRunListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Flow pipeline runs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "list_pipeline_runs", params)
			if err != nil {
				return err
			}
			if options.JSONOutput {
				_, _ = fmt.Fprintln(streams.Out, result)
				return nil
			}
			return printPipelineRunList(streams.Out, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.PipelineID, "pipeline-id", "", "Flow pipeline ID")
	flags.IntVar(&options.Page, "page", 0, "page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "page size")
	flags.Int64Var(&options.StartTime, "start-time", 0, "run start time as Unix milliseconds")
	flags.Int64Var(&options.EndTime, "end-time", 0, "run end time as Unix milliseconds")
	flags.StringVar(&options.Status, "status", "", "run status, e.g. FAIL, SUCCESS, or RUNNING")
	flags.IntVar(&options.TriggerMode, "trigger-mode", 0, "trigger mode: 1 manual, 2 scheduled, 3 code push, 5 pipeline, 6 webhook")
	flags.BoolVar(&options.JSONOutput, "json", false, "print raw JSON")
	return command
}

func (o pipelineRunListOptions) params() (map[string]any, error) {
	params := map[string]any{
		"pipelineId": strings.TrimSpace(o.PipelineID),
	}
	if params["pipelineId"] == "" {
		return nil, fmt.Errorf("pipeline-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "status", o.Status)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	if o.StartTime > 0 {
		params["startTime"] = o.StartTime
	}
	if o.EndTime > 0 {
		params["endTime"] = o.EndTime
	}
	if o.TriggerMode > 0 {
		params["triggerMode"] = o.TriggerMode
	}
	return params, nil
}

func printPipelineRunList(out anyWriter, raw string) error {
	rows, ok := pipelineRunRowsFromJSONForPrint(raw)
	if !ok {
		_, _ = fmt.Fprintln(out, raw)
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tSTATUS\tRESULT\tSTART\tEND\tTRIGGER")
	for _, row := range rows {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", row.ID, row.Status, row.Result, row.Start, row.End, row.Trigger)
	}
	return writer.Flush()
}

type pipelineRunRow struct {
	ID      string
	Status  string
	Result  string
	Start   string
	End     string
	Trigger string
}

func pipelineRunRowsFromJSON(raw string) []pipelineRunRow {
	rows, _ := pipelineRunRowsFromJSONForPrint(raw)
	return rows
}

func pipelineRunRowsFromJSONForPrint(raw string) ([]pipelineRunRow, bool) {
	items, ok := rowsFromJSONWithPresence(raw)
	if !ok {
		return nil, false
	}
	rows := make([]pipelineRunRow, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, pipelineRunRow{
			ID:      firstStringValue(m, "pipelineRunId", "runId", "id"),
			Status:  firstStringValue(m, "status", "runStatus", "state"),
			Result:  firstStringValue(m, "result", "runResult", "executeResult"),
			Start:   firstStringValue(m, "startTime", "startedAt", "gmtStarted", "createdAt"),
			End:     firstStringValue(m, "endTime", "endTme", "finishedAt", "gmtFinished", "updatedAt"),
			Trigger: firstStringValue(m, "triggerMode", "triggerType", "triggerUser", "creator"),
		})
	}
	return rows, true
}
