package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version info populated by main.go from ldflags
var (
	buildVersion = "dev"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

// SetVersionInfo is called by main.go to populate version information
func SetVersionInfo(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date
}

// GetVersion returns the semantic version string
func GetVersion() string {
	return buildVersion
}

// GetCommit returns the git commit hash
func GetCommit() string {
	return buildCommit
}

// GetDate returns the build date
func GetDate() string {
	return buildDate
}

// IsDebugBuild returns true if this is a dev build (not set via ldflags)
func IsDebugBuild() bool {
	return buildVersion == "dev"
}

func init() {

}

var rootCmd = &cobra.Command{
	Use:   "bdcli",
	Short: "CLI for managing BetterDiscord",
	Long:  `A cross-platform CLI for installing, updating, and managing BetterDiscord.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("You should probably use a subcommand")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
