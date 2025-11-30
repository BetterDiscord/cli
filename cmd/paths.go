package cmd

import (
    "fmt"

    "github.com/betterdiscord/cli/internal/discord"
    "github.com/betterdiscord/cli/internal/models"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(pathsCmd)
}

var pathsCmd = &cobra.Command{
    Use:   "paths",
    Short: "Show suggested Discord install paths",
    Long:  "Displays the suggested core installation path per Discord channel detected on this system.",
    Run: func(cmd *cobra.Command, args []string) {
        channels := []models.DiscordChannel{models.Stable, models.PTB, models.Canary}
        for _, ch := range channels {
            p := discord.GetSuggestedPath(ch)
            name := ch.Name()
            if p == "" {
                fmt.Printf("%s: (none detected)\n", name)
            } else {
                fmt.Printf("%s: %s\n", name, p)
            }
        }
    },
}
