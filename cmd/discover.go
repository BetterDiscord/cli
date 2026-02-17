package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	discoverCmd.AddCommand(discoverInstallsCmd)
	discoverCmd.AddCommand(discoverPathsCmd)
	discoverCmd.AddCommand(discoverAddonsCmd)
	rootCmd.AddCommand(discoverCmd)
}

var discoverCmd = &cobra.Command{
	Use:     "discover",
	// Aliases: []string{"info", "list"},
	Short:   "Discover Discord installations and related data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return discoverInstallsCmd.RunE(cmd, args)
	},
}

var discoverInstallsCmd = &cobra.Command{
	Use:   "installs",
	Short: "Show detected Discord installations",
	Long:  "Lists detected Discord installations by channel, showing path, version, install type, and BetterDiscord status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		installs := discord.GetAllInstalls()
		if len(installs) == 0 {
			fmt.Println("No Discord installations detected.")
			return nil
		}

		channels := []models.DiscordChannel{models.Stable, models.PTB, models.Canary}
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "CHANNEL\tVERSION\tTYPE\tBD INJECTED\tPATH")

		for _, ch := range channels {
			arr := installs[ch]
			for _, inst := range arr {
				typeLabel := "native"
				if inst.IsFlatpak {
					typeLabel = "flatpak"
				} else if inst.IsSnap {
					typeLabel = "snap"
				}
				bdStatus := "no"
				if inst.IsInjected() {
					bdStatus = "yes"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", ch.Name(), inst.Version, typeLabel, bdStatus, inst.CorePath)
			}
		}

		return tw.Flush()
	},
}

var discoverPathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Show suggested install paths per channel",
	RunE: func(cmd *cobra.Command, args []string) error {
		channels := []models.DiscordChannel{models.Stable, models.PTB, models.Canary}
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "CHANNEL\tSUGGESTED PATH")
		for _, ch := range channels {
			p := discord.GetSuggestedPath(ch)
			if p == "" {
				p = "(none detected)"
			}
			fmt.Fprintf(tw, "%s\t%s\n", ch.Name(), p)
		}
		return tw.Flush()
	},
}

var discoverAddonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "Summarize installed plugins and themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		plugins, err := betterdiscord.ListAddons(betterdiscord.AddonPlugin)
		if err != nil {
			return err
		}
		themes, err := betterdiscord.ListAddons(betterdiscord.AddonTheme)
		if err != nil {
			return err
		}

		fmt.Printf("Plugins: %d installed\n", len(plugins))
		for _, p := range plugins {
			fmt.Printf("  - %s (%s)\n", p.FullFilename, p.Path)
		}
		fmt.Printf("Themes: %d installed\n", len(themes))
		for _, t := range themes {
			fmt.Printf("  - %s (%s)\n", t.FullFilename, t.Path)
		}
		return nil
	},
}
