package cmd

import "github.com/spf13/cobra"

// 根命令直接显示用法，子命令在用法前显示说明。
const helpTemplate = `{{if .HasParent}}{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}{{if not .HasParent}}{{range .Commands}}{{if eq .Name "version"}}
Version Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}{{end}}{{end}}`

func init() {
	mainCommand.SetHelpTemplate(helpTemplate)

	// 隐藏版本子命令的帮助选项，旧用法直接跳转到总帮助。
	commandVersion.Flags().BoolP("help", "h", false, "")
	_ = commandVersion.Flags().MarkHidden("help")
	versionHelp := commandVersion.HelpFunc()
	commandVersion.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("help") {
			mainCommand.InitDefaultHelpFlag()
			mainCommand.HelpFunc()(mainCommand, nil)
			return
		}
		versionHelp(cmd, args)
	})
}
