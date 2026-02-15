package cmd

import (
	"fmt"
	"path"

	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	uninstallCmd.Flags().StringP("path", "p", "", "Path to a Discord installation")
	uninstallCmd.Flags().StringP("channel", "c", "stable", "Discord release channel (stable|ptb|canary)")
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls BetterDiscord from your Discord",
	Long:  "Uninstall BetterDiscord by specifying either --path to a Discord install or --channel to auto-detect (default: stable).",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathFlag, _ := cmd.Flags().GetString("path")
		channelFlag, _ := cmd.Flags().GetString("channel")

		pathProvided := pathFlag != ""
		channelProvided := cmd.Flags().Changed("channel")

		if pathProvided && channelProvided {
			return fmt.Errorf("--path and --channel are mutually exclusive")
		}

		var install *discord.DiscordInstall

		if pathProvided {
			install = discord.ResolvePath(pathFlag)
			if install == nil {
				return fmt.Errorf("could not find a valid Discord installation at %s", pathFlag)
			}
		} else {
			channel := models.ParseChannel(channelFlag)
			resolvedPath := discord.GetSuggestedPath(channel)
			install = discord.ResolvePath(resolvedPath)
			if install == nil {
				return fmt.Errorf("could not find a valid %s installation to uninstall", channelFlag)
			}
		}

		if err := install.UninstallBD(); err != nil {
			return fmt.Errorf("uninstallation failed: %w", err)
		}

		fmt.Printf("✅ BetterDiscord uninstalled from %s\n", path.Dir(install.CorePath))
		return nil
	},
}
