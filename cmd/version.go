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
  Long:  "Shows all the necessary version information for bug reports.",
  Run: func(cmd *cobra.Command, args []string) {
	fmt.Println(debug.ReadBuildInfo())
  },
}