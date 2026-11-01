package command

import (
	"context"
	"fmt"
	"os"

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
		Short: "Run scenarios",
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
	// fail-fast
	scenarios := make([]scenarioFile, 0, len(opts.scenarioFiles))
	for _, path := range opts.scenarioFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read scenario file %s: %w", path, err)
		}
		scenarios = append(scenarios, scenarioFile{
			path: path,
			data: data,
		})
	}

	if opts.configFile != "" {
		// apply config
		if err := runApply(ctx, cli, applyOptions{filePath: opts.configFile}); err != nil {
			return err
		}
	}

	// run scenarios
	for _, sc := range scenarios {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("run is interrupted: %w", err)
		}

		// TODO: parsing + send sc.data into usernetd daemon
		fmt.Fprintf(cli.Out(), "Scenario %q (%d bytes) is completed succesfully", sc.path, len(sc.data))
	}

	return nil
}
