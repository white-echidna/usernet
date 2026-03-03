package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/white-echidna/usernet/internal/app"
	"github.com/white-echidna/usernet/internal/topology"
)

var (
	rootCmd = &cobra.Command{
		Use:   "usernet",
		Short: "A user space network emulator",
		Run:   run,
	}
	file string
)

func init() {
	rootCmd.Flags().StringVarP(&file, "file", "f", "topology.json", "topology file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	topology, err := topology.Parse(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := app.New(topology)

	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
