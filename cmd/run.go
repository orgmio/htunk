package cmd

import (
	"github.com/orgmio/htunk/internal/proxy"
	"github.com/spf13/cobra"
)

func init() {
	var options proxy.Options
	commandRun := &cobra.Command{
		Use:   "run",
		Short: "Run the proxy service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return proxy.Run(cmd.Context(), options)
		},
	}
	commandRun.Flags().StringVarP(&options.ConfigPath, "config", "c", "config.soml", "set configuration file path")
	commandRun.Flags().BoolVar(&options.DisableColor, "disable-color", false, "disable color output")
	mainCommand.AddCommand(commandRun)
}
