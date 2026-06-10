package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoCompareOptions struct {
	OrganizationID string
	RepositoryID   string
	SourceType     string
	TargetType     string
	Straight       string
}

func newYunxiaoRepoCompareCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoCompareOptions
	command := &cobra.Command{
		Use:     "compare <from> <to>",
		Aliases: []string{"diff"},
		Short:   "compare CodeUp repository refs as JSON",
		Example: `  # Compare two branches
  yunxiao repo compare main..feature --repository-id group/repo

  # Compare two commits
  yunxiao repo compare abc123 def456 --repository-id group/repo`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params(args[0], args[1])
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "compare", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "CodeUp repository ID or full path")
	flags.StringVar(&options.SourceType, "source-type", "", "source ref type, e.g. branch, tag, or commit")
	flags.StringVar(&options.TargetType, "target-type", "", "target ref type, e.g. branch, tag, or commit")
	flags.StringVar(&options.Straight, "straight", "", "whether to compare directly without merge base, e.g. true or false")
	return command
}

func (o repoCompareOptions) params(from, to string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"from":         strings.TrimSpace(from),
		"to":           strings.TrimSpace(to),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["from"] == "" {
		return nil, fmt.Errorf("from is required")
	}
	if params["to"] == "" {
		return nil, fmt.Errorf("to is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "sourceType", o.SourceType)
	setCLIStringParam(params, "targetType", o.TargetType)
	setCLIStringParam(params, "straight", o.Straight)
	return params, nil
}
