package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "usernet COMMAND [ARG...]",
		Short:         "A user space network emulator",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetIn(cli.In())
	cmd.SetOut(cli.Out())
	cmd.SetErr(cli.Err())

	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%w\nTry '%s --help'", err, c.CommandPath())
	})

	cmd.AddCommand(
		NewApplyCommand(cli),
		NewRunCommand(cli),
	)

	return cmd
}
