package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/spf13/cobra"
)

func init() {
	// Store parent command with subcommands
	storeCmd.AddCommand(storeSearchCmd)
	storeCmd.AddCommand(storeShowCmd)
	storeCmd.AddCommand(storePluginsCmd)
	storeCmd.AddCommand(storeThemesCmd)

	// Plugins subcommands
	storePluginsCmd.AddCommand(storePluginsSearchCmd)
	storePluginsCmd.AddCommand(storePluginsShowCmd)

	// Themes subcommands
	storeThemesCmd.AddCommand(storeThemesSearchCmd)
	storeThemesCmd.AddCommand(storeThemesShowCmd)

	// Register to root
	rootCmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Browse and search the BetterDiscord store",
	Long:  "Search and view addon information from the BetterDiscord store.",
}

// ==================== Store (search all addons) ====================

var storeSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for addons",
	Long:  "Search all addons in the BetterDiscord store.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		addons, err := betterdiscord.FetchAddonsOfType("")
		if err != nil {
			return err
		}

		results := betterdiscord.SearchAddons(addons, query)
		if len(results) == 0 {
			fmt.Println("No addons found matching that query.")
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tTYPE\tVERSION\tAUTHOR\tDOWNLOADS")
		for _, addon := range results {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%d\n", addon.ID, addon.Name, addon.Type, addon.Version, addon.Author.DisplayName, addon.Downloads)
		}
		return tw.Flush()
	},
}

var storeShowCmd = &cobra.Command{
	Use:   "show <id|name>",
	Short: "Show addon details",
	Long:  "Show detailed information about an addon from the BetterDiscord store.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		addon, err := betterdiscord.FetchAddonFromStore(identifier)
		if err != nil {
			return err
		}

		betterdiscord.LogAddonInfo(addon)
		return nil
	},
}

// ==================== Plugins ====================

var storePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Search and view plugins",
	Long:  "Browse and search plugins from the BetterDiscord store.",
}

var storePluginsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for plugins",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		addons, err := betterdiscord.FetchAddonsOfType("plugins")
		if err != nil {
			return err
		}

		results := betterdiscord.SearchAddons(addons, query)
		if len(results) == 0 {
			fmt.Println("No plugins found matching that query.")
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tVERSION\tAUTHOR\tDOWNLOADS")
		for _, addon := range results {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\n", addon.ID, addon.Name, addon.Version, addon.Author.DisplayName, addon.Downloads)
		}
		return tw.Flush()
	},
}

var storePluginsShowCmd = &cobra.Command{
	Use:   "show <id|name>",
	Short: "Show plugin details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		addon, err := betterdiscord.FetchAddonFromStore(identifier)
		if err != nil {
			return err
		}

		betterdiscord.LogAddonInfo(addon)
		return nil
	},
}

// ==================== Themes ====================

var storeThemesCmd = &cobra.Command{
	Use:   "themes",
	Short: "Search and view themes",
	Long:  "Browse and search themes from the BetterDiscord store.",
}

var storeThemesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for themes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		addons, err := betterdiscord.FetchAddonsOfType("themes")
		if err != nil {
			return err
		}

		results := betterdiscord.SearchAddons(addons, query)
		if len(results) == 0 {
			fmt.Println("No themes found matching that query.")
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tVERSION\tAUTHOR\tDOWNLOADS")
		for _, addon := range results {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\n", addon.ID, addon.Name, addon.Version, addon.Author.DisplayName, addon.Downloads)
		}
		return tw.Flush()
	},
}

var storeThemesShowCmd = &cobra.Command{
	Use:   "show <id|name>",
	Short: "Show theme details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		addon, err := betterdiscord.FetchAddonFromStore(identifier)
		if err != nil {
			return err
		}

		betterdiscord.LogAddonInfo(addon)
		return nil
	},
}
