package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineRunViewOptions struct {
	OrganizationID string
	PipelineID     string
	IncludeJobs    bool
	Category       string
}

func newYunxiaoPipelineRunViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := pipelineRunViewOptions{
		IncludeJobs: true,
	}
	command := &cobra.Command{
		Use:     "view <pipeline-run-id>",
		Aliases: []string{"overview"},
		Short:   "view a Flow pipeline run overview as JSON",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params(args[0])
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_pipeline_run_overview", params)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.PipelineID, "pipeline-id", "", "Flow pipeline ID")
	flags.BoolVar(&options.IncludeJobs, "include-jobs", true, "include pipeline jobs by category")
	flags.StringVar(&options.Category, "category", "", "pipeline job category; defaults to DEPLOY")
	return command
}

func (o pipelineRunViewOptions) params(pipelineRunID string) (map[string]any, error) {
	params := map[string]any{
		"pipelineId":    strings.TrimSpace(o.PipelineID),
		"pipelineRunId": strings.TrimSpace(pipelineRunID),
	}
	if params["pipelineId"] == "" {
		return nil, fmt.Errorf("pipeline-id is required")
	}
	if params["pipelineRunId"] == "" {
		return nil, fmt.Errorf("pipeline-run-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "category", o.Category)
	params["includeJobs"] = o.IncludeJobs
	return params, nil
}
