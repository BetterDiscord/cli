package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	utils "betterdiscord/cli/utils"
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
		var targetExe = "s"
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

		// Make BD directories
		if err := os.MkdirAll(utils.Data, 0755); err != nil {
			fmt.Println("Could not create BetterDiscord folder")
			return
		}

		if err := os.MkdirAll(utils.Plugins, 0755); err != nil {
			fmt.Println("Could not create plugins folder")
			return
		}

		if err := os.MkdirAll(utils.Themes, 0755); err != nil {
			fmt.Println("Could not create theme folder")
			return
		}

		// Get download URL from GitHub API
		var apiData, err = utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
		if err != nil {
			fmt.Println("Could not get API response")
			fmt.Println(err)
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
			fmt.Println("Could not download asar")
			return
		}

		// Inject shim loader
		var discordPath = utils.DiscordPath(releaseChannel)
		var appPath = path.Join(discordPath, "app")
		if err := os.MkdirAll(appPath, 0755); err != nil {
			fmt.Println("Could not create app folder")
			return
		}

		var pkgPath = path.Join(appPath, "package.json")
		var indPath = path.Join(appPath, "index.js")
		const adpString = `{"name":"betterdiscord","main":"index.js"}`
		const ddcString = `{"name":"discord_desktop_core","version":"0.0.0","private":"true","main":"index.js"}`
		var pkgString = adpString
		var indString = `require("` + asarPath + `");`
		indString = strings.ReplaceAll(indString, `\`, "/")

		if runtime.GOOS == "linux" {
			pkgString = ddcString
			indString = indString + `\nmodule.exports = require("./core.asar")`
		}

		if err := os.WriteFile(pkgPath, []byte(pkgString), 0755); err != nil {
			fmt.Println("Could not write package.json")
			return
		}

		if err := os.WriteFile(indPath, []byte(indString), 0755); err != nil {
			fmt.Println("Could not write index.js")
			return
		}

		// Rename Asar
		if runtime.GOOS != "linux" {
			var appAsarPath = path.Join(discordPath, "app.asar")
			var disAsarPath = path.Join(discordPath, "discord.asar")
			var appAsarExists = utils.Exists(appAsarPath)
			var disAsarExists = utils.Exists(disAsarPath)

			// If neither exists, something is really wrong
			if !appAsarExists && !disAsarExists {
				fmt.Println("Discord installation corrupt")
				return
			}

			// If both exist, get rid of previously renamed asar as outdated
			if appAsarExists && disAsarExists {
				if err := os.Remove(disAsarPath); err != nil {
					fmt.Println("Could not delete discord.asar, is Discord running?")
					return
				}
			}

			// If the app asar exists, rename it. This will also handle cases from
			// above where both existed so do not make this an else or !disAsarExists
			if appAsarExists {
				if err := os.Rename(appAsarPath, disAsarPath); err != nil {
					fmt.Println("Could not rename app.asar, is Discord running?")
					return
				}
			}
		}

		// Launch Discord if we killed it
		if len(exe) > 0 {
			var cmd = exec.Command(exe)
			cmd.Start()
		}
	},
}
