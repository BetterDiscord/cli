package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
)

func init() {
	installCmd.Flags().StringP("path", "p", "", "Path to a Discord installation")
	installCmd.Flags().StringP("channel", "c", "stable", "Discord release channel (stable|ptb|canary)")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{"update"},
	Short: "Installs BetterDiscord to your Discord",
	Long:  "Install BetterDiscord by specifying either --path to a Discord install or --channel to auto-detect (default: stable).",
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
			corePath := discord.GetSuggestedPath(channel)
			install = discord.ResolvePath(corePath)
			if install == nil {
				return fmt.Errorf("could not find a valid %s installation to install to", channelFlag)
			}
		}

		if err := install.InstallBD(); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		fmt.Printf("✅ BetterDiscord installed to %s\n", path.Dir(install.CorePath))
		return nil
	},
}
