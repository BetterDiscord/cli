package cmd

import (
	"fmt"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	// Parent command: themes
	themesCmd.AddCommand(themesListCmd)
	themesCmd.AddCommand(themesInfoCmd)
	themesCmd.AddCommand(themesInstallCmd)
	themesCmd.AddCommand(themesRemoveCmd)
	themesCmd.AddCommand(themesUpdateCmd)
	rootCmd.AddCommand(themesCmd)
}

var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "Manage BetterDiscord themes",
	Long:  "List, install, remove, and update BetterDiscord themes.",
}

var themesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := betterdiscord.ListAddons(betterdiscord.AddonTheme)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			output.Println("📭 No themes installed.")
			return nil
		}

		tw := output.NewTableWriter()
		fmt.Fprintln(tw, "NAME\tVERSION\tAUTHOR\tSIZE (KB)\tMODIFIED")
		for _, item := range items {
			name := item.Meta.Name
			if name == "" {
				name = item.BaseName
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%.1f\t%s\n", name, item.Meta.Version, item.Meta.Author, float64(item.Size)/1024.0, item.Modified.Format(output.DateTimeFormat))
		}
		return tw.Flush()
	},
}

var themesInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed information about an installed theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		items, err := betterdiscord.ListAddons(betterdiscord.AddonTheme)
		if err != nil {
			return err
		}

		for _, item := range items {
			// Match by filename or meta name
			if item.BaseName == name || item.FullFilename == name || item.Meta.Name == name {
				betterdiscord.LogLocalAddonInfo(&item)
				return nil
			}
		}

		output.Printf("❌ Theme '%s' not found.\n", name)
		return nil
	},
}

var themesInstallCmd = &cobra.Command{
	Use:   "install <name|id|url>",
	Short: "Install a theme by name, ID, or direct URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		resolved, err := betterdiscord.InstallAddon(betterdiscord.AddonTheme, identifier)
		if err != nil {
			return err
		}
		output.Printf("✅ Theme installed at %s\n", resolved.URL)
		return nil
	},
}

var themesRemoveCmd = &cobra.Command{
	Use:     "remove <name|id>",
	Aliases: []string{"uninstall"},
	Short:   "Remove an installed theme",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		if err := betterdiscord.RemoveAddon(betterdiscord.AddonTheme, identifier); err != nil {
			return err
		}
		output.Printf("✅ Theme removed: %s\n", identifier)
		return nil
	},
}

var themesUpdateCmd = &cobra.Command{
	Use:   "update <name|id|url>",
	Short: "Update a theme by name, ID, or URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		resolved, err := betterdiscord.UpdateAddon(betterdiscord.AddonTheme, identifier)
		if err != nil {
			return err
		}
		output.Printf("✅ Theme updated at %s\n", resolved.URL)
		return nil
	},
}
