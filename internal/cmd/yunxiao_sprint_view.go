package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type sprintViewOptions struct {
	OrganizationID string
	ProjectID      string
	Categories     string
	Subject        string
	Status         string
	AssignedTo     string
	Creator        string
	SampleLimit    int
	SampleLimitSet bool
}

func newYunxiaoSprintViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options sprintViewOptions
	command := &cobra.Command{
		Use:     "view <sprint-id>",
		Aliases: []string{"overview"},
		Short:   "view a Projex sprint overview as JSON",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			options.SampleLimitSet = cmd.Flags().Changed("sample-limit")
			params, err := options.params(args[0])
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_sprint_overview", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.ProjectID, "project-id", "", "Projex project ID")
	flags.StringVar(&options.Categories, "categories", "", "comma-separated work item categories; defaults to Task,Bug")
	flags.StringVar(&options.Subject, "subject", "", "subject/title keyword")
	flags.StringVar(&options.Status, "status", "", "comma-separated status IDs")
	flags.StringVar(&options.AssignedTo, "assigned-to", "", "comma-separated assignee user IDs")
	flags.StringVar(&options.Creator, "creator", "", "comma-separated creator user IDs")
	flags.IntVar(&options.SampleLimit, "sample-limit", 0, "work item samples returned per category")
	return command
}

func (o sprintViewOptions) params(sprintID string) (map[string]any, error) {
	params := map[string]any{
		"projectId": strings.TrimSpace(o.ProjectID),
		"sprintId":  strings.TrimSpace(sprintID),
	}
	if params["projectId"] == "" {
		return nil, fmt.Errorf("project-id is required")
	}
	if params["sprintId"] == "" {
		return nil, fmt.Errorf("sprint-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "categories", o.Categories)
	setCLIStringParam(params, "subject", o.Subject)
	setCLIStringParam(params, "status", o.Status)
	setCLIStringParam(params, "assignedTo", o.AssignedTo)
	setCLIStringParam(params, "creator", o.Creator)
	if o.SampleLimitSet {
		params["sampleLimit"] = o.SampleLimit
	}
	return params, nil
}
