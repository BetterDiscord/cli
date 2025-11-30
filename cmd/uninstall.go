package cmd

import (
	"fmt"

	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:       "uninstall <channel>",
	Short:     "Uninstalls BetterDiscord from your Discord",
	Long:      "This can uninstall BetterDiscord to multiple versions and paths of Discord at once. Options for channel are: stable, canary, ptb",
	ValidArgs: []string{"canary", "stable", "ptb"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var releaseChannel = args[0]
		var corePath = discord.GetSuggestedPath(models.ParseChannel(releaseChannel))
		var install = discord.ResolvePath(corePath)

		if install == nil {
			fmt.Printf("❌ Could not find a valid %s installation to uninstall.\n", releaseChannel)
			return
		}

		if err := install.UninstallBD(); err != nil {
			fmt.Printf("❌ Uninstallation failed: %s\n", err.Error())
			return
		}

		fmt.Printf("✅ BetterDiscord uninstalled from %s\n", corePath)
	},
}
