// cmd 包实现 htunk 的命令行界面。
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

// Execute 根据命令行参数执行对应的命令。
func Execute() {
	configureHelp()
	if err := mainCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
