package cmd

import (
	"errors"
	"fmt"

	"github.com/orgmio/htunk/internal/config"
	"github.com/orgmio/htunk/internal/log"
	"github.com/orgmio/htunk/internal/proxy"
	"github.com/spf13/cobra"
)

func init() {
	var configPath string
	var disableColor bool
	commandRun := &cobra.Command{
		Use:   "run",
		Short: "Run the proxy service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				cmd.SilenceErrors = true
				log.NewConsole(cmd.ErrOrStderr(), disableColor).ErrorContext(cmd.Context(), err.Error())
				return err
			}
			logger, closeLog, err := log.Open(cfg.Log, cmd.ErrOrStderr(), disableColor)
			if err != nil {
				cmd.SilenceErrors = true
				log.NewConsole(cmd.ErrOrStderr(), disableColor).ErrorContext(cmd.Context(), err.Error())
				return err
			}

			runErr := proxy.Run(cmd.Context(), cfg)
			if runErr != nil {
				err = fmt.Errorf("proxy: %w", runErr)
				logger.ErrorContext(cmd.Context(), err.Error())
			}
			if closeErr := closeLog(); closeErr != nil {
				return errors.Join(err, fmt.Errorf("close log file: %w", closeErr))
			}
			cmd.SilenceErrors = true
			return err
		},
	}
	mainCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "config.soml", "set configuration file path")
	mainCommand.PersistentFlags().BoolVar(&disableColor, "disable-color", false, "disable color output")
	mainCommand.AddCommand(commandRun)
}
