package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/white-echidna/usernet/control_plane/cli/command"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	usernetCli := command.NewUsernetCli()
	rootCmd := command.NewRootCommand(usernetCli)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(usernetCli.Err(), "Error: %v\n", err)
		os.Exit(1)
	}
}
