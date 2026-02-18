package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/betterdiscord/cli/internal/output"
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

// isURL checks if a string is a valid URL
func isURL(input string) bool {
	parsed, err := url.Parse(input)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

var silent bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Suppress non-error output")
}

var rootCmd = &cobra.Command{
	Use:   "bdcli",
	Short: "CLI for managing BetterDiscord",
	Long:  `A cross-platform CLI for installing, updating, and managing BetterDiscord.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if silent || isSilentEnvEnabled() {
			output.SetWriters(io.Discard, nil)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
}

func isSilentEnvEnabled() bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("BDCLI_SILENT")))
	return value != "" && value != "0" && value != "false" && value != "no"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(output.ErrorWriter(), err)
		os.Exit(1)
	}
}
