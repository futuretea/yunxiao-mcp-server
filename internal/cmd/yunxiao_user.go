package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newYunxiaoUserCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		Short:   "work with Yunxiao users",
	}
	command.AddCommand(newYunxiaoUserWhoamiCommand(streams, cfgFile, v))
	return command
}

func newYunxiaoUserWhoamiCommand(streams IOStreams, cfgFile *string, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:   "whoami",
		Short: "show the current authenticated user as JSON",
		Example: `  # Show current user
  yunxiao user whoami`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadYunxiaoCLIConfig(cmd, *cfgFile, v)
			if err != nil {
				return err
			}
			result, err := callYunxiaoTool(cmd, cfg, "get_current_user", map[string]any{})
			if err != nil {
				return err
			}
			printCLIJSON(streams.Out, result)
			return nil
		},
	}
	return command
}
