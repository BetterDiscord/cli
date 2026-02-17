package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/output"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about BetterDiscord installation",
	Long:  "Displays detailed information about the BetterDiscord installation, including version, commit, branch, build, and installation paths.",
	RunE: func(cmd *cobra.Command, args []string) error {
		bdinstall := betterdiscord.GetInstallation()

		if !bdinstall.IsAsarInstalled() {
			return fmt.Errorf("BetterDiscord does not appear to be installed, try running 'bdcli install' first")
		}

		output.Printf("📦 BetterDiscord Information:\n\n")

		bdinstall.LogBuildinfo()

		return nil
	},
}
