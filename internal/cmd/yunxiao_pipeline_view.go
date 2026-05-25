package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pipelineViewOptions struct {
	OrganizationID string
	IncludeRuns    bool
	RunLimit       int
}

func newYunxiaoPipelineViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := pipelineViewOptions{
		IncludeRuns: true,
	}
	command := &cobra.Command{
		Use:     "view <pipeline-id>",
		Aliases: []string{"overview"},
		Short:   "view a Flow pipeline overview as JSON",
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
			result, err := callYunxiaoTool(cmd, cfg, "get_pipeline_overview", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeRuns, "include-runs", true, "include recent run history")
	flags.IntVar(&options.RunLimit, "run-limit", 0, "max recent runs returned")
	return command
}

func (o pipelineViewOptions) params(pipelineID string) (map[string]any, error) {
	params := map[string]any{
		"pipelineId": strings.TrimSpace(pipelineID),
	}
	if params["pipelineId"] == "" {
		return nil, fmt.Errorf("pipeline-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	params["includeRuns"] = o.IncludeRuns
	if o.RunLimit > 0 {
		params["runLimit"] = o.RunLimit
	}
	return params, nil
}
