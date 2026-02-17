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

		buildinfo, err := bdinstall.ReadBuildinfo()
		if err != nil {
			return err
		}

		output.Printf("📦 BetterDiscord Information:\n\n")

		output.Printf("   Build Information:\n")
		output.Printf("     🔹 Version: %s\n", output.FormatVersion(buildinfo.Version))
		output.Printf("     🔹 Commit:  %s\n", buildinfo.Commit)
		output.Printf("     🔹 Branch:  %s\n", buildinfo.Branch)
		output.Printf("     🔹 Mode:    %s\n\n", buildinfo.Mode)

		output.Printf("   Installation Paths:\n")
		output.Printf("     📁 Base:    %s\n", bdinstall.Root())
		output.Printf("     ⚙️  Data:    %s\n", bdinstall.Data())
		output.Printf("     🔌 Plugins: %s\n", bdinstall.Plugins())
		output.Printf("     🎨 Themes:  %s\n", bdinstall.Themes())

		return nil
	},
}
