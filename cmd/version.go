package cmd

import (
	"runtime/debug"

	"github.com/betterdiscord/cli/internal/output"
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
			output.Println(debug.ReadBuildInfo())
		} else {
			output.Printf("📦 BetterDiscord CLI %s\n", GetVersion())
			output.Printf("🔖 Commit: %s\n", GetCommit())
			output.Printf("🕒 Built:  %s\n", GetDate())
		}
	},
}
