package command

import (
	"context"

	"github.com/spf13/cobra"
)

type runOptions struct {
	configFile    string
	scenarioFiles []string
}

func NewRunCommand(cli Cli) *cobra.Command {
	var opts runOptions

	cmd := &cobra.Command{
		Use:   "run [-f CONFIG_FILE] SCENARIO_FILE [SCENARIO_FILE...]",
		Short: "Run specified scenario",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.scenarioFiles = args
			return runScenarios(cmd.Context(), cli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFile, "file", "f", "", "Optional path to YAML config")
	_ = cmd.MarkFlagFilename("file", "yaml", "yml")

	return cmd
}

type scenarioFile struct {
	path string
	data []byte
}

func runScenarios(ctx context.Context, cli Cli, opts runOptions) error {
	panic("not implemented")
}
