package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/white-echidna/usernet/control_plane/cli"
)

func NewRootCommand(usernetCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "usernet COMMAND [ARG...]",
		Short:         "A user space network emulator",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetIn(usernetCli.In())
	cmd.SetOut(usernetCli.Out())
	cmd.SetErr(usernetCli.Err())

	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%w\nTry '%s --help'", err, c.CommandPath())
	})

	cmd.AddCommand(
		NewApplyCommand(usernetCli),
		NewRunCommand(usernetCli),
	)

	return cmd
}
