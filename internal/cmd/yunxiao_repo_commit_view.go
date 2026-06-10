package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoCommitViewOptions struct {
	OrganizationID   string
	RepositoryID     string
	IncludeStatuses  bool
	IncludeCheckRuns bool
	StatusLimit      int
	CheckRunLimit    int
}

func newYunxiaoRepoCommitViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := repoCommitViewOptions{
		IncludeStatuses:  true,
		IncludeCheckRuns: true,
	}
	command := &cobra.Command{
		Use:     "view <sha>",
		Aliases: []string{"detail"},
		Short:   "view a CodeUp commit overview as JSON",
		Example: `  # View commit by SHA
  yunxiao repo commit view abc123def --repository-id group/repo`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			params, err := options.params(args[0])
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_commit_overview", params)
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
	flags.BoolVar(&options.IncludeStatuses, "include-statuses", true, "include commit status records")
	flags.BoolVar(&options.IncludeCheckRuns, "include-check-runs", true, "include check run records")
	flags.IntVar(&options.StatusLimit, "status-limit", 0, "max commit statuses returned")
	flags.IntVar(&options.CheckRunLimit, "check-run-limit", 0, "max check runs returned")
	return command
}

func (o repoCommitViewOptions) params(sha string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"sha":          strings.TrimSpace(sha),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["sha"] == "" {
		return nil, fmt.Errorf("sha is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	params["includeStatuses"] = o.IncludeStatuses
	params["includeCheckRuns"] = o.IncludeCheckRuns
	if o.StatusLimit > 0 {
		params["statusLimit"] = o.StatusLimit
	}
	if o.CheckRunLimit > 0 {
		params["checkRunLimit"] = o.CheckRunLimit
	}
	return params, nil
}
