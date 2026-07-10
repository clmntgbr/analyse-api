package main

import (
	"fmt"
	cliCommand "go-api/cmd/cli/command"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "Go CLI - cmd commands",
		Long:  "Go CLI provides commands for cmd tasks",
	}

	rootCmd.AddCommand(
		cliCommand.NewMigrateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
