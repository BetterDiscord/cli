package cmd

import (
    "fmt"
    "os"
    "path"
    "strings"

    "github.com/spf13/cobra"

    utils "betterdiscord/cli/utils"
)

func init() {
    rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
  Use:   "install",
  Short: "Installs BetterDiscord to your Discord",
  Long:  "This can install BetterDiscord to multiple versions and paths of Discord at once.",
  Run: func(cmd *cobra.Command, args []string) {
    var dir, _ = os.UserConfigDir();
    var betterdiscordDir = path.Join(dir, "BetterDiscord")
    var testDir = path.Join(dir, "ekfjkwnfwek")
    var _, bdErr = os.ReadDir(betterdiscordDir)
    var _, tErr = os.ReadDir(testDir)
    var localAppData = os.Getenv("LOCALAPPDATA")
    var discordDir = path.Join(localAppData, "Discord")
    var dFiles, _ = os.ReadDir(discordDir)
    var appDir = ""
    for _, file := range dFiles {
        if !file.IsDir() || !strings.HasPrefix(file.Name(), "app-") {
            continue;
        }
        if (file.Name() > appDir) {
            appDir = file.Name()
        }
    }
    fmt.Println(appDir)
    fmt.Printf("bd exists %t | ekfjkwnfwek exists %t | %s \n", bdErr == nil, tErr == nil, localAppData);

    var shouldKill, _ = cmd.Flags().GetBool("kill-discord")
    fmt.Printf("Should we kill %t\n", shouldKill)
    if shouldKill {
        if err := utils.KillProcess("Discord.exe"); err != nil {
            fmt.Println("Could not kill Discord")
            return
        }
    }

    fmt.Println(utils.DiscordPath("stable"))
    
    // fmt.Println(err)
    // var data, err = utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest");
	// if err != nil {
	// 	fmt.Println("Could not get API response")
	// 	fmt.Println(err)
	// 	return
	// }
	// var index = 0;
	// for i, asset := range data.Assets {
	// 	if asset.Name == "betterdiscord.asar" {
	// 		index = i;
	// 		break;
	// 	}
	// }
	// fmt.Println(data.Assets[index].URL)
  },
};