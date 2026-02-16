package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Shows version information for BetterDiscord CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		// Show clean output for production builds, debug info for dev builds
		if IsDebugBuild() {
			fmt.Println(debug.ReadBuildInfo())
		} else {
			fmt.Printf("BetterDiscord CLI %s\n", GetVersion())
			fmt.Printf("Commit: %s\n", GetCommit())
			fmt.Printf("Built:  %s\n", GetDate())
		}
	},
}

