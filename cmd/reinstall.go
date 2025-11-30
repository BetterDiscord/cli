package cmd

import (
    "fmt"

    "github.com/betterdiscord/cli/internal/discord"
    "github.com/betterdiscord/cli/internal/models"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(reinstallCmd)
}

var reinstallCmd = &cobra.Command{
    Use:       "reinstall <channel>",
    Short:     "Uninstall and then reinstall BetterDiscord",
    Long:      "Performs an uninstall followed by an install of BetterDiscord for the specified Discord channel.",
    ValidArgs: []string{"canary", "stable", "ptb"},
    Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
    Run: func(cmd *cobra.Command, args []string) {
        releaseChannel := args[0]
        corePath := discord.GetSuggestedPath(models.ParseChannel(releaseChannel))
        inst := discord.ResolvePath(corePath)
        if inst == nil {
            fmt.Printf("❌ Could not find a valid %s installation to reinstall.\n", releaseChannel)
            return
        }

        if err := inst.UninstallBD(); err != nil {
            fmt.Printf("❌ Uninstall failed: %s\n", err.Error())
            return
        }

        if err := inst.InstallBD(); err != nil {
            fmt.Printf("❌ Install failed: %s\n", err.Error())
            return
        }

        fmt.Println("✅ BetterDiscord reinstalled successfully")
    },
}
