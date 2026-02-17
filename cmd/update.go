package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

func init() {
	updateCmd.Flags().BoolP("check", "c", false, "Only check for updates, don't install")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update BetterDiscord to the latest version",
	Long:  "Download and install the latest version of BetterDiscord.",
	RunE: func(cmd *cobra.Command, args []string) error {
		bdinstall := betterdiscord.GetInstallation()

		if !bdinstall.IsAsarInstalled() {
			return fmt.Errorf("BetterDiscord does not appear to be installed, run 'bdcli install' first")
		}

		checkFlag, _ := cmd.Flags().GetBool("check")

		// Get current version
		buildinfo, err := bdinstall.ReadBuildinfo()
		if err != nil {
			return fmt.Errorf("failed to read current BetterDiscord version: %w", err)
		}

		currentVersion := buildinfo.Version
		fmt.Printf("📦 Current version: v%s\n", currentVersion)

		// Get latest release from GitHub
		release, err := utils.DownloadJSON[models.GitHubRelease]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		latestVersion := release.TagName
		fmt.Printf("🌐 Latest version:  %s\n\n", latestVersion)

		// Check if update is needed
		if compareVersions(currentVersion, latestVersion) >= 0 {
			fmt.Printf("✅ You are already on the latest version!\n")
			return nil
		}

		fmt.Printf("🎉 New version available!\n\n")

		if checkFlag {
			fmt.Println("Run 'bdcli update' to install the update")
			return nil
		}

		// Download the latest version
		fmt.Println("📥 Downloading update...")
		if err := bdinstall.Download(); err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}

		fmt.Printf("✅ Successfully updated to %s\n", latestVersion)
		return nil
	},
}

// compareVersions compares two semantic versions (e.g., "1.0.156" vs "1.0.157")
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Strip 'v' prefix if present
	v1 = "" + v1
	v2 = "" + v2
	if len(v1) > 0 && v1[0] == 'v' {
		v1 = v1[1:]
	}
	if len(v2) > 0 && v2[0] == 'v' {
		v2 = v2[1:]
	}

	// Parse into version parts
	parts1 := splitVersion(v1)
	parts2 := splitVersion(v2)

	// Compare each part
	maxLen := max(len(parts2), len(parts1))

	for i := range maxLen {
		var p1, p2 int

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		} else if p1 > p2 {
			return 1
		}
	}

	return 0
}

// splitVersion splits a version string into parts (e.g., "1.0.156" -> ["1", "0", "156"])
func splitVersion(v string) []string {
	var parts []string
	var current string
	for i := 0; i < len(v); i++ {
		if v[i] == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else if v[i] >= '0' && v[i] <= '9' {
			current += string(v[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
