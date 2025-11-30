package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:       "install <channel>",
	Short:     "Installs BetterDiscord to your Discord",
	Long:      "This can install BetterDiscord to multiple versions and paths of Discord at once. Options for channel are: stable, canary, ptb",
	ValidArgs: []string{"canary", "stable", "ptb"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var releaseChannel = args[0]
		var corePath = discord.GetSuggestedPath(models.ParseChannel(releaseChannel))
		var install = discord.ResolvePath(corePath)

		if install == nil {
			fmt.Printf("❌ Could not find a valid %s installation to install to.\n", releaseChannel)
			return
		}

		if err := install.InstallBD(); err != nil {
			fmt.Printf("❌ Installation failed: %s\n", err.Error())
			return
		}

		fmt.Printf("✅ BetterDiscord installed to %s\n", path.Dir(install.GetPath()))
	},
}
