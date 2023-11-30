package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/utils"
)

func strinlist(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func init() {
	rootCmd.AddCommand(installCmd)
}
func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}

var ValidArgs = []string{"canary", "stable", "ptb"}
var installCmd = &cobra.Command{
	Use:       "install <channel>",
	Short:     "Installs BetterDiscord to your Discord",
	Long:      "This can install BetterDiscord to multiple versions and paths of Discord at once. Options for channel are: stable, canary, ptb",
	ValidArgs: ValidArgs,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var releaseChannel = args[0]
		if !strinlist(releaseChannel, ValidArgs) {
			die("invalid arguments given. valid arguments are: " + strings.Join(ValidArgs, ", "))
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
			targetExe = "Discord.exe"
		}

		// Kill Discord if it's running
		var exe = utils.GetProcessExe(targetExe)
		if len(exe) > 0 {
			if err := utils.KillProcess(targetExe); err != nil {
				die("Could not kill Discord")
				return
			}
		}

		// Make BD directories
		if err := os.MkdirAll(utils.Data, 0755); err != nil {
			die("Could not create BetterDiscord folder")
			return
		}

		if err := os.MkdirAll(utils.Plugins, 0755); err != nil {
			die("Could not create plugins folder")
			return
		}

		if err := os.MkdirAll(utils.Themes, 0755); err != nil {
			die("Could not create theme folder")
			return
		}

		// Get download URL from GitHub API
		var apiData, err = utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
		if err != nil {
			fmt.Println("Could not get API response")
			die(err.Error())
			return
		}

		var index = 0
		for i, asset := range apiData.Assets {
			if asset.Name == "betterdiscord.asar" {
				index = i
				break
			}
		}

		var downloadUrl = apiData.Assets[index].URL

		// Download asar into the BD folder
		var asarPath = path.Join(utils.Data, "betterdiscord.asar")
		err = utils.DownloadFile(downloadUrl, asarPath)
		if err != nil {
			die("Could not download asar")
			return
		}

		// Inject shim loader
		var corePath = utils.DiscordPath(releaseChannel)

		var indString = `require("` + asarPath + `");`
		indString = strings.ReplaceAll(indString, `\`, "/")
		indString = indString + "\nmodule.exports = require(\"./core.asar\");"

		if err := os.WriteFile(path.Join(corePath, "index.js"), []byte(indString), 0755); err != nil {
			die("Could not write index.js in discord_desktop_core!")
			return
		}

		// Launch Discord if we killed it
		if len(exe) > 0 {
			var cmd = exec.Command(exe)
			_ = cmd.Start()
		}
	},
}
