package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repoFileViewOptions struct {
	OrganizationID string
	RepositoryID   string
	Ref            string
	Since          string
}

func newYunxiaoRepoFileViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	var options repoFileViewOptions
	command := &cobra.Command{
		Use:     "view <path>",
		Aliases: []string{"content"},
		Short:   "view a CodeUp repository file as JSON",
		Example: `  # View file content
  yunxiao repo file view README.md --repository-id group/repo

  # View file on a specific branch
  yunxiao repo file view src/main.go --repository-id group/repo --ref-name develop`,
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
			result, err := callYunxiaoTool(cmd, cfg, "get_file_blobs", params)
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
	flags.StringVar(&options.Ref, "ref", "", "branch, tag, or commit SHA")
	flags.StringVar(&options.Since, "since", "", "ISO 8601 timestamp to retrieve the file as it existed then")
	return command
}

func (o repoFileViewOptions) params(path string) (map[string]any, error) {
	params := map[string]any{
		"repositoryId": strings.TrimSpace(o.RepositoryID),
		"filePath":     strings.TrimSpace(path),
		"ref":          strings.TrimSpace(o.Ref),
	}
	if params["repositoryId"] == "" {
		return nil, fmt.Errorf("repository-id is required")
	}
	if params["filePath"] == "" {
		return nil, fmt.Errorf("path is required")
	}
	if params["ref"] == "" {
		return nil, fmt.Errorf("ref is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "since", o.Since)
	return params, nil
}
