package cmd

import (
    "fmt"

    "github.com/betterdiscord/cli/internal/discord"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List detected Discord installations",
    Long:  "Scans common locations and lists detected Discord installations grouped by channel.",
    Run: func(cmd *cobra.Command, args []string) {
        installs := discord.GetAllInstalls()
        if len(installs) == 0 {
            fmt.Println("No Discord installations detected.")
            return
        }
        for channel, arr := range installs {
            if len(arr) == 0 {
                continue
            }
            fmt.Printf("%s:\n", channel.Name())
            for _, inst := range arr {
                fmt.Printf("  - %s (version %s)\n", inst.GetPath(), discord.GetVersion(inst.GetPath()))
            }
        }
    },
}
