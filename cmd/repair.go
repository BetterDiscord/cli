package cmd

import (
    "fmt"

    "github.com/betterdiscord/cli/internal/discord"
    "github.com/betterdiscord/cli/internal/models"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(repairCmd)
}

var repairCmd = &cobra.Command{
    Use:       "repair <channel>",
    Short:     "Repairs the BetterDiscord installation",
    Long:      "Attempts to repair the BetterDiscord setup for the specified Discord channel (e.g., disables problematic plugins).",
    ValidArgs: []string{"canary", "stable", "ptb"},
    Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
    Run: func(cmd *cobra.Command, args []string) {
        releaseChannel := args[0]
        corePath := discord.GetSuggestedPath(models.ParseChannel(releaseChannel))
        inst := discord.ResolvePath(corePath)
        if inst == nil {
            fmt.Printf("❌ Could not find a valid %s installation to repair.\n", releaseChannel)
            return
        }
        if err := inst.RepairBD(); err != nil {
            fmt.Printf("❌ Repair failed: %s\n", err.Error())
            return
        }
        fmt.Println("✅ Repair completed successfully")
    },
}
