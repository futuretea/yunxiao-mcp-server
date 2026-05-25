package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type taskViewOptions struct {
	OrganizationID     string
	IncludeActivities  bool
	IncludeRelations   bool
	RelationTypes      string
	IncludeAttachments bool
	IncludeComments    bool
	Page               int
	PerPage            int
}

func newYunxiaoTaskViewCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	options := taskViewOptions{
		IncludeActivities:  true,
		IncludeRelations:   true,
		RelationTypes:      "ASSOCIATED,SUB",
		IncludeAttachments: true,
		IncludeComments:    true,
	}
	command := &cobra.Command{
		Use:     "view <workitem-id>",
		Aliases: []string{"detail"},
		Short:   "view a Projex task detail as JSON",
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
			result, err := callYunxiaoTool(cmd, cfg, "get_project_workitem_detail", params)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(streams.Out, result)
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.OrganizationID, "organization-id", "", "Yunxiao organization ID; defaults when the token belongs to one organization")
	flags.BoolVar(&options.IncludeActivities, "include-activities", true, "include task activity history")
	flags.BoolVar(&options.IncludeRelations, "include-relations", true, "include task relation records")
	flags.StringVar(&options.RelationTypes, "relation-types", "ASSOCIATED,SUB", "comma-separated relation types when relations are included")
	flags.BoolVar(&options.IncludeAttachments, "include-attachments", true, "include task attachments")
	flags.BoolVar(&options.IncludeComments, "include-comments", true, "include task comments")
	flags.IntVar(&options.Page, "page", 0, "comments page number")
	flags.IntVar(&options.PerPage, "per-page", 0, "comments page size")
	return command
}

func (o taskViewOptions) params(workitemID string) (map[string]any, error) {
	params := map[string]any{
		"workitemId":         strings.TrimSpace(workitemID),
		"includeActivities":  o.IncludeActivities,
		"includeRelations":   o.IncludeRelations,
		"includeAttachments": o.IncludeAttachments,
		"includeComments":    o.IncludeComments,
	}
	if params["workitemId"] == "" {
		return nil, fmt.Errorf("workitem-id is required")
	}
	setCLIStringParam(params, "organizationId", o.OrganizationID)
	setCLIStringParam(params, "relationTypes", o.RelationTypes)
	if o.Page > 0 {
		params["page"] = o.Page
	}
	if o.PerPage > 0 {
		params["perPage"] = o.PerPage
	}
	return params, nil
}
