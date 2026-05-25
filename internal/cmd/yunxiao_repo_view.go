package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoViewOptions struct {
	OrganizationID       string
	IncludeBranches      bool
	IncludeCommits       bool
	IncludeMergeRequests bool
	Ref                  string
	BranchLimit          int
	CommitLimit          int
	MergeRequestLimit    int
	MergeRequestState    string
}

func newYunxiaoRepoViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := repoViewOptions{
		IncludeBranches:      true,
		IncludeCommits:       true,
		IncludeMergeRequests: true,
	}
	command := &cobra.Command{
		Use:     "view <repository-id>",
		Aliases: []string{"overview"},
		Short:   "view a CodeUp repository overview as JSON",
		Example: `  # View repository by full path
  yunxiao repo view group/repo

  # View with explicit organization
  yunxiao repo view group/repo --organization-id org-abc`,
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
			result, err := callYunxiaoTool(cmd, cfg, "get_repository_overview", params)
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeBranches, "include-branches", true, "include branch records")
	flags.BoolVar(&options.IncludeCommits, "include-commits", true, "include recent commit records")
	flags.BoolVar(&options.IncludeMergeRequests, "include-merge-requests", true, "include merge request records")
	flags.StringVar(&options.Ref, "ref", "", "branch, tag, or commit SHA for commit listing")
	flags.IntVar(&options.BranchLimit, "branch-limit", 0, "max branches returned")
	flags.IntVar(&options.CommitLimit, "commit-limit", 0, "max commits returned")
	flags.IntVar(&options.MergeRequestLimit, "merge-request-limit", 0, "max merge requests returned")
	flags.StringVar(&options.MergeRequestState, "merge-request-state", "", "merge request state filter, e.g. opened, merged, or closed")
	return command
}

func (o repoViewOptions) params(repositoryID string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(repositoryID),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	params["includeBranches"] = o.IncludeBranches
	params["includeCommits"] = o.IncludeCommits
	params["includeMergeRequests"] = o.IncludeMergeRequests
	setCLIStringParam(params, "refName", o.Ref)
	if o.BranchLimit > 0 {
		params["branchLimit"] = o.BranchLimit
	}
	if o.CommitLimit > 0 {
		params["commitLimit"] = o.CommitLimit
	}
	if o.MergeRequestLimit > 0 {
		params["mrLimit"] = o.MergeRequestLimit
	}
	setCLIStringParam(params, "mrState", o.MergeRequestState)
	return params, nil
}
