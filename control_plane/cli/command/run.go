package command

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/white-echidna/usernet/control_plane/cli"
)

type runOptions struct {
	filePath  string
	scenarios []string
}

func NewRunCommand(usernetCli cli.Cli) *cobra.Command {
	var opts runOptions

	cmd := &cobra.Command{
		Use:   "run -f FILENAME SCENARIO [SCENARIO...]",
		Short: "Run specified scenario",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.scenarios = args
			return runScenarios(cmd.Context(), usernetCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.filePath, "file", "f", "", "Path to configuration file")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagFilename("file", "yaml", "yml")

	return cmd
}

func runScenarios(ctx context.Context, usernetCli cli.Cli, opts runOptions) error {
	panic("not implemented")
}
