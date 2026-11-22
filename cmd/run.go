package cmd

import (
	"github.com/orgmio/htunk/internal/log"
	"github.com/orgmio/htunk/internal/proxy"
	"github.com/spf13/cobra"
)

func init() {
	var options proxy.Options
	var disableColor bool
	commandRun := &cobra.Command{
		Use:   "run",
		Short: "Run the proxy service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := log.NewConsole(cmd.ErrOrStderr(), disableColor)
			if err := proxy.Run(cmd.Context(), options); err != nil {
				cmd.SilenceErrors = true
				logger.ErrorContext(cmd.Context(), "proxy: "+err.Error())
				return err
			}
			return nil
		},
	}
	mainCommand.PersistentFlags().StringVarP(&options.ConfigPath, "config", "c", "config.soml", "set configuration file path")
	mainCommand.PersistentFlags().BoolVar(&disableColor, "disable-color", false, "disable color output")
	mainCommand.AddCommand(commandRun)
}
