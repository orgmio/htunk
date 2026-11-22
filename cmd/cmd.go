// Package cmd implements the htunk command-line interface.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var mainCommand = &cobra.Command{
	Use:          "htunk",
	Short:        "A general-purpose proxy platform",
	SilenceUsage: true,
}

// Execute runs the command selected by the command-line arguments.
func Execute() {
	if err := mainCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
