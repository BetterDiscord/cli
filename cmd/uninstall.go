package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	utils "github.com/betterdiscord/cli/internal/utils"
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
		var corePath = utils.DiscordPath(releaseChannel)
		var indString = "module.exports = require(\"./core.asar\");"

		if err := os.WriteFile(path.Join(corePath, "index.js"), []byte(indString), 0755); err != nil {
			fmt.Println("Could not write index.js in discord_desktop_core!")
			return
		}

		var targetExe = ""
		switch releaseChannel {
		case "stable":
			targetExe = "Discord.exe"
			break
		case "canary":
			targetExe = "DiscordCanary.exe"
			break
		case "ptb":
			targetExe = "DiscordPTB.exe"
			break
		default:
			targetExe = ""
		}

		// Kill Discord if it's running
		var exe = utils.GetProcessExe(targetExe)
		if len(exe) > 0 {
			if err := utils.KillProcess(targetExe); err != nil {
				fmt.Println("Could not kill Discord")
				return
			}
		}

		// Launch Discord if we killed it
		if len(exe) > 0 {
			var cmd = exec.Command(exe)
			cmd.Start()
		}
	},
}
