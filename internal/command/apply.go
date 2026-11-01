package command

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type applyOptions struct {
	filePath string
}

func NewApplyCommand(cli Cli) *cobra.Command {
	var opts applyOptions

	cmd := &cobra.Command{
		Use:   "apply -f FILENAME",
		Short: "Load YAML config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply(cmd.Context(), cli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.filePath, "file", "f", "", "Path to YAML config")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagFilename("file", "yaml", "yml")

	return cmd
}

func runApply(ctx context.Context, cli Cli, opts applyOptions) error {
	data, err := os.ReadFile(opts.filePath)
	if err != nil {
		return fmt.Errorf("cannot read config file %s: %w", opts.filePath, err)
	}

	// TODO: parse config; send config and ctx into usernetd
	fmt.Fprintf(cli.Out(), "Config %q (%d bytes) is applied successfully", opts.filePath, len(data))

	return nil
}
