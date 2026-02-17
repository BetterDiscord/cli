package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/spf13/cobra"
)

func init() {
	// Parent command: plugins
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsCmd.AddCommand(pluginsInfoCmd)
	pluginsCmd.AddCommand(pluginsInstallCmd)
	pluginsCmd.AddCommand(pluginsRemoveCmd)
	pluginsCmd.AddCommand(pluginsUpdateCmd)
	rootCmd.AddCommand(pluginsCmd)
}

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage BetterDiscord plugins",
	Long:  "List, install, remove, and update BetterDiscord plugins.",
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := betterdiscord.ListAddons(betterdiscord.AddonPlugin)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			fmt.Println("No plugins installed.")
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "NAME\tVERSION\tAUTHOR\tSIZE (KB)\tMODIFIED")
		for _, item := range items {
			name := item.Meta.Name
			if name == "" {
				name = item.BaseName
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%.1f\t%s\n", name, item.Meta.Version, item.Meta.Author, float64(item.Size)/1024.0, item.Modified.Format("2006-01-02 15:04"))
		}
		return tw.Flush()
	},
}

var pluginsInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show detailed information about an installed plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		items, err := betterdiscord.ListAddons(betterdiscord.AddonPlugin)
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

		fmt.Printf("Plugin '%s' not found\n", name)
		return nil
	},
}

var pluginsInstallCmd = &cobra.Command{
	Use:   "install <name|id|url>",
	Short: "Install a plugin by name, ID, or direct URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		resolved, err := betterdiscord.InstallAddon(betterdiscord.AddonPlugin, identifier)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Plugin installed at %s\n", resolved.URL)
		return nil
	},
}

var pluginsRemoveCmd = &cobra.Command{
	Use:     "remove <name|id>",
	Aliases: []string{"uninstall"},
	Short:   "Remove an installed plugin",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		if err := betterdiscord.RemoveAddon(betterdiscord.AddonPlugin, identifier); err != nil {
			return err
		}
		fmt.Printf("Removed plugin %s\n", identifier)
		return nil
	},
}

var pluginsUpdateCmd = &cobra.Command{
	Use:   "update <name|id|url>",
	Short: "Update a plugin by name, ID, or URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		resolved, err := betterdiscord.UpdateAddon(betterdiscord.AddonPlugin, identifier)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Plugin updated at %s\n", resolved.URL)
		return nil
	},
}
