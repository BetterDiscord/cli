package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/betterdiscord"
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

		fmt.Printf("📦 BetterDiscord Information:\n\n")

		fmt.Printf("   Build Information:\n")
		fmt.Printf("     🔹 Version: v%s\n", buildinfo.Version)
		fmt.Printf("     🔹 Commit:  %s\n", buildinfo.Commit)
		fmt.Printf("     🔹 Branch:  %s\n", buildinfo.Branch)
		fmt.Printf("     🔹 Mode:    %s\n\n", buildinfo.Mode)

		fmt.Printf("   Installation Paths:\n")
		fmt.Printf("     📁 Base:    %s\n", bdinstall.Root())
		fmt.Printf("     ⚙️  Data:    %s\n", bdinstall.Data())
		fmt.Printf("     🔌 Plugins: %s\n", bdinstall.Plugins())
		fmt.Printf("     🎨 Themes:  %s\n", bdinstall.Themes())

		return nil
	},
}
