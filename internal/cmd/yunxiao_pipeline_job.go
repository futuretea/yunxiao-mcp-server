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
	command.AddCommand(newYunxiaoPipelineJobLogCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoPipelineJobListCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineJobListOptions
	command := &cobra.Command{
		Use:   "list",
		Short: "list Flow pipeline jobs by category",
		Example: `  # List deploy jobs
  yunxiao pipeline job list --pipeline-id pipeline-123 --category DEPLOY

  # Output as JSON
  yunxiao pipeline job list --pipeline-id pipeline-123 --category DEPLOY --json`,
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
				printCLIJSON(streams.Out, result)
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
		_, _ = fmt.Fprintln(out, "No results found.")
		return nil
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, boldTableHeader("IDENTIFIER\tNAME\tCATEGORY\tSTATUS"))
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

type pipelineJobLogOptions struct {
	OrganizationID string
	PipelineID     string
	PipelineRunID  string
	JobID          string
}

func newYunxiaoPipelineJobLogCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options pipelineJobLogOptions
	command := &cobra.Command{
		Use:   "log",
		Short: "get execution log for a pipeline job",
		Example: `  # Get job log
  yunxiao pipeline job log --pipeline-id pipeline-123 --run-id run-456 --job-id job-789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params()
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_pipeline_job_run_log", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.PipelineID, "pipeline-id", "", "Flow pipeline ID")
	flags.StringVar(&options.PipelineRunID, "run-id", "", "pipeline run ID")
	flags.StringVar(&options.JobID, "job-id", "", "pipeline job ID")
	return command
}

func (o pipelineJobLogOptions) params() (map[string]any, error) {
	params := map[string]any{
		"pipelineId":    strings.TrimSpace(o.PipelineID),
		"pipelineRunId": strings.TrimSpace(o.PipelineRunID),
		"jobId":         strings.TrimSpace(o.JobID),
	}
	if params["pipelineId"] == "" {
		return nil, fmt.Errorf("pipeline-id is required")
	}
	if params["pipelineRunId"] == "" {
		return nil, fmt.Errorf("run-id is required")
	}
	if params["jobId"] == "" {
		return nil, fmt.Errorf("job-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	return params, nil
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
