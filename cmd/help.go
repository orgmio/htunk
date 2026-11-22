package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// 根命令直接显示用法，子命令在用法前显示说明。
const helpTemplate = `{{if .HasParent}}{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}{{if not .HasParent}}{{range .Commands}}{{if eq .Name "version"}}
Version Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}{{end}}{{end}}`

func configureHelp() {
	mainCommand.InitDefaultHelpFlag()
	mainCommand.InitDefaultHelpCmd()
	mainCommand.SetUsageTemplate(strings.ReplaceAll(mainCommand.UsageTemplate(),
		"{{.CommandPath}} [command] --help", "{{.CommandPath}} help [command]"))
	mainCommand.SetHelpTemplate(helpTemplate)

	helpCommand, _, _ := mainCommand.Find([]string{"help"})
	helpCommand.Run = nil
	helpCommand.RunE = func(cmd *cobra.Command, args []string) error {
		target, remaining, err := cmd.Root().Find(args)
		if err != nil || len(remaining) != 0 {
			return fmt.Errorf("unknown help topic %q", strings.Join(args, " "))
		}

		target.InitDefaultHelpFlag()
		target.InitDefaultVersionFlag()
		return target.Help()
	}

	// 隐藏版本子命令的帮助选项，旧用法直接跳转到总帮助。
	commandVersion.Flags().BoolP("help", "h", false, "")
	_ = commandVersion.Flags().MarkHidden("help")
	versionHelp := commandVersion.HelpFunc()
	commandVersion.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("help") {
			_ = cmd.Root().Help()
			return
		}
		versionHelp(cmd, args)
	})
}
