package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoChangeRequestViewOptions struct {
	OrganizationID   string
	RepositoryID     string
	IncludePatchSets bool
	IncludeComments  bool
	CommentState     string
	CommentResolved  bool
}

func newYunxiaoRepoChangeRequestViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := repoChangeRequestViewOptions{
		IncludePatchSets: true,
		IncludeComments:  true,
	}
	command := &cobra.Command{
		Use:     "view <local-id>",
		Aliases: []string{"detail", "overview"},
		Short:   "view a CodeUp change request overview as JSON",
		Example: `  # View change request by local ID
  yunxiao repo change-request view 42 --repository-id group/repo`,
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
			result, err := callYunxiaoTool(cmd, cfg, "get_change_request_overview", params)
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
	flags.BoolVar(&options.IncludePatchSets, "include-patch-sets", true, "include patch set records")
	flags.BoolVar(&options.IncludeComments, "include-comments", true, "include comment records")
	flags.StringVar(&options.CommentState, "comment-state", "", "comment state filter, e.g. OPENED or RESOLVED")
	flags.BoolVar(&options.CommentResolved, "comment-resolved", false, "include resolved comments")
	return command
}

func (o repoChangeRequestViewOptions) params(localID string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"localId":      strings.TrimSpace(localID),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["localId"] == "" {
		return nil, fmt.Errorf("local-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	params["includePatchSets"] = o.IncludePatchSets
	params["includeComments"] = o.IncludeComments
	setCLIStringParam(params, "commentState", o.CommentState)
	params["commentResolved"] = o.CommentResolved
	return params, nil
}
