package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/betterdiscord/cli/utils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var ValidArgs = []string{"stable", "canary", "ptb"}
var installCmd = &cobra.Command{
	Use:       "install [" + strings.Join(ValidArgs, ", ") + "]",
	Short:     "Installs BetterDiscord to your Discord",
	Long:      "This can install BetterDiscord to multiple versions and paths of Discord at once. Options for channel are: stable, canary, ptb",
	ValidArgs: ValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var releaseChannel = ""
		if len(args) < 1 {
			releaseChannel = "stable"
		} else {
			releaseChannel = args[0]
		}

		if !slices.Contains(ValidArgs, releaseChannel) {
			log.Fatal("invalid arguments given. valid arguments are: " + strings.Join(ValidArgs, ", "))
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
				log.Fatal("Could not kill Discord")
				return
			}
		}

		// Make BD directories
		if err := os.MkdirAll(utils.Data, 0755); err != nil {
			log.Fatal("Could not create BetterDiscord folder")
		}

		if err := os.MkdirAll(utils.Plugins, 0755); err != nil {
			log.Fatal("Could not create plugins folder")
		}

		if err := os.MkdirAll(utils.Themes, 0755); err != nil {
			log.Fatal("Could not create theme folder")
		}

		// Get download URL from GitHub API
		var apiData, err = utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
		if err != nil {
			fmt.Println("Could not get API response")
			log.Fatal(err.Error())
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
			log.Fatal("Could not download asar")
		}

		// Inject shim loader
		var corePath = utils.DiscordPath(releaseChannel)

		var indString = `require("` + asarPath + `");`
		indString = strings.ReplaceAll(indString, `\`, "/")
		indString = indString + "\nmodule.exports = require(\"./core.asar\");"

		if err := os.WriteFile(path.Join(corePath, "index.js"), []byte(indString), 0755); err != nil {
			log.Fatal("Could not write index.js in discord_desktop_core!")
		}

		// Launch Discord if we killed it
		if len(exe) > 0 {
			var cmd = exec.Command(exe)
			_ = cmd.Start()
		}
	},
}
