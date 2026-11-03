package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/white-echidna/usernet/control_plane/cli"
)

type applyOptions struct {
	configPath string
}

func NewApplyCommand(usernetCli cli.Cli) *cobra.Command {
	var opts applyOptions

	cmd := &cobra.Command{
		Use:   "apply -f FILENAME",
		Short: "Apply configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply(cmd.Context(), usernetCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configPath, "file", "f", "", "Path to configuration file")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagFilename("file", "yaml", "yml")

	return cmd
}

func runApply(ctx context.Context, usernetCli cli.Cli, opts applyOptions) error {
	panic("not implemented")
}
