package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoBranchViewOptions struct {
	OrganizationID       string
	RepositoryID         string
	IncludeCommits       bool
	IncludeMergeRequests bool
	CommitLimit          int
	MergeRequestLimit    int
	MergeRequestState    string
}

func newYunxiaoRepoBranchViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := repoBranchViewOptions{
		IncludeCommits:       true,
		IncludeMergeRequests: true,
	}
	command := &cobra.Command{
		Use:     "view <branch-name>",
		Aliases: []string{"overview", "detail"},
		Short:   "view a CodeUp branch overview as JSON",
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
			result, err := callYunxiaoTool(cmd, cfg, "get_branch_overview", params)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.StringVar(&options.RepositoryID, "repository-id", "", "CodeUp repository ID or full path")
	flags.BoolVar(&options.IncludeCommits, "include-commits", true, "include recent commit records")
	flags.BoolVar(&options.IncludeMergeRequests, "include-merge-requests", true, "include merge request records")
	flags.IntVar(&options.CommitLimit, "commit-limit", 0, "max commits returned")
	flags.IntVar(&options.MergeRequestLimit, "merge-request-limit", 0, "max merge requests returned")
	flags.StringVar(&options.MergeRequestState, "merge-request-state", "", "merge request state filter, e.g. opened, merged, or closed")
	return command
}

func (o repoBranchViewOptions) params(branchName string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"branchName":   strings.TrimSpace(branchName),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["branchName"] == "" {
		return nil, fmt.Errorf("branch-name is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	params["includeCommits"] = o.IncludeCommits
	params["includeMergeRequests"] = o.IncludeMergeRequests
	if o.CommitLimit > 0 {
		params["commitLimit"] = o.CommitLimit
	}
	if o.MergeRequestLimit > 0 {
		params["mrLimit"] = o.MergeRequestLimit
	}
	setCLIStringParam(params, "mrState", o.MergeRequestState)
	return params, nil
}
