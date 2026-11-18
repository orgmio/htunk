package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var nameOnly bool

var commandVersion = &cobra.Command{
	Use:   "version",
	Short: "Print current version of htunk",
	Args:  cobra.NoArgs,
	Run: func(*cobra.Command, []string) {
		if nameOnly {
			fmt.Println(Version)
			return
		}
		PrintVersion()
	},
}

func init() {
	commandVersion.Flags().BoolVarP(&nameOnly, "name", "n", false, "print version name only")
	mainCommand.AddCommand(commandVersion)
}

func PrintVersion() {
	fmt.Println("htunk v" + Version)
	fmt.Println()
	settings := make(map[string]string)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			settings[s.Key] = s.Value
		}
	}
	fmt.Println("Environment: " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH)
	fmt.Println("Revision: " + settings["vcs.revision"])
	cgo := "disabled"
	if settings["CGO_ENABLED"] == "1" {
		cgo = "enabled"
	}
	fmt.Println("CGO_ENABLED: " + cgo)
}
