package cmd

import (
	"fmt"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	initThemesCmd()
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
		existing := betterdiscord.FindAddon(betterdiscord.AddonTheme, name)
        if existing == nil {
            output.Printf("❌ Theme '%s' not found.\n", name)
            return nil
        }

        betterdiscord.LogLocalAddonInfo(existing)
		return nil
	},
}

var themesInstallCmd = &cobra.Command{
	Use:   "install <name|id|url>",
	Short: "Install a theme by name, ID, or direct URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		// Check if not a URL and already installed
		if !utils.IsURL(identifier) {
			if existing := betterdiscord.FindAddon(betterdiscord.AddonTheme, identifier); existing != nil {
				name := existing.Meta.Name
				if name == "" {
					name = existing.BaseName
				}
				output.Printf("⚠️ Theme '%s' is already installed.\n", name)
				output.Println("💡 To update the theme, use: bdcli themes update <name|id|url>")
				return nil
			}
		}
		resolved, err := betterdiscord.InstallAddon(betterdiscord.AddonTheme, identifier)
		if err != nil {
			return err
		}
		output.Printf("✅ Theme installed at %s\n", resolved.Path)
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
		// Check if addon exists before attempting removal
		existing := betterdiscord.FindAddon(betterdiscord.AddonTheme, identifier)
		if existing == nil {
			output.Printf("❌ Theme '%s' is not installed.\n", identifier)
			return nil
		}
		if err := betterdiscord.RemoveAddon(betterdiscord.AddonTheme, identifier); err != nil {
			return err
		}
		name := existing.Meta.Name
		if name == "" {
			name = existing.BaseName
		}
		output.Printf("✅ Theme removed: %s\n", name)
		return nil
	},
}

var themesUpdateCmd = &cobra.Command{
	Use:   "update <name|id|url>",
	Short: "Update a theme by name, ID, or URL",
	Args: func(cmd *cobra.Command, args []string) error {
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag {
			if len(args) > 0 {
				return fmt.Errorf("cannot specify addon identifier when using --all flag")
			}
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")
		allFlag, _ := cmd.Flags().GetBool("all")

		// Handle --all flag
		if allFlag {
			return updateAllThemes(checkOnly)
		}

		identifier := args[0]

		// For non-URL identifiers, check if update is available
		if !utils.IsURL(identifier) {
			existing := betterdiscord.FindAddon(betterdiscord.AddonTheme, identifier)
			if existing == nil {
				output.Printf("❌ Theme '%s' is not installed.\n", identifier)
				return nil
			}

			// Try to fetch from store to check version
			store, err := betterdiscord.FetchAddonFromStore(identifier)
			if err == nil && store != nil {
				localVersion := existing.Meta.Version
				storeVersion := store.Version

				if localVersion == storeVersion {
					localName := existing.Meta.Name
					if localName == "" {
						localName = existing.BaseName
					}
					output.Printf("✅ Theme '%s' is already up to date (v%s)\n", localName, localVersion)
					return nil
				}

				localName := existing.Meta.Name
				if localName == "" {
					localName = existing.BaseName
				}

				if checkOnly {
					output.Printf("📦 Update available for '%s'\n", localName)
					output.Printf("   Current: v%s → Available: v%s\n", localVersion, storeVersion)
					output.Println("💡 To install the update, use: bdcli themes update <name|id|url> (without --check)")
					return nil
				}
			}
		}

		if checkOnly {
			output.Println("⚠️  Cannot check version when using direct URL")
			return nil
		}

		resolved, err := betterdiscord.UpdateAddon(betterdiscord.AddonTheme, identifier)
		if err != nil {
			return err
		}
		output.Printf("✅ Theme updated at %s\n", resolved.Path)
		return nil
	},
}

func updateAllThemes(checkOnly bool) error {
	items, err := betterdiscord.ListAddons(betterdiscord.AddonTheme)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		output.Println("📭 No themes installed.")
		return nil
	}

	var toUpdate []struct {
		entry        betterdiscord.AddonEntry
		localVersion string
		storeVersion string
		name         string
	}

	output.Println("🔍 Checking for theme updates...")

	// Check each theme for updates
	for _, item := range items {
		name := item.Meta.Name
		if name == "" {
			name = item.BaseName
		}

		// Try to find in store by name or ID
		identifier := name
		if identifier == "" {
			identifier = item.BaseName
		}

		store, err := betterdiscord.FetchAddonFromStore(identifier)
		if err != nil {
			// Theme not in store, skip
			continue
		}

		localVersion := item.Meta.Version
		storeVersion := store.Version

		if localVersion != storeVersion {
			toUpdate = append(toUpdate, struct {
				entry        betterdiscord.AddonEntry
				localVersion string
				storeVersion string
				name         string
			}{
				entry:        item,
				localVersion: localVersion,
				storeVersion: storeVersion,
				name:         name,
			})
		}
	}

	if len(toUpdate) == 0 {
		output.Println("✅ All themes are up to date!")
		return nil
	}

	output.Printf("\n📦 Found %d theme(s) with available updates:\n\n", len(toUpdate))

	// Show what would be updated
	for _, item := range toUpdate {
		output.Printf("  • %s: v%s → v%s\n", item.name, item.localVersion, item.storeVersion)
	}
	output.Println()

	if checkOnly {
		output.Println("💡 To install these updates, use: bdcli themes update --all (without --check)")
		return nil
	}

	// Perform updates
	updated := 0
	failed := 0

	for _, item := range toUpdate {
		identifier := item.name
		if identifier == "" {
			identifier = item.entry.BaseName
		}

		_, err := betterdiscord.UpdateAddon(betterdiscord.AddonTheme, identifier)
		if err != nil {
			output.Printf("❌ Failed to update %s: %v\n", item.name, err)
			failed++
		} else {
			output.Printf("✅ Updated %s to v%s\n", item.name, item.storeVersion)
			updated++
		}
	}

	output.Printf("\n📊 Summary: %d updated, %d failed\n", updated, failed)
	return nil
}

func initThemesCmd() {
	// Parent command: themes
	themesCmd.AddCommand(themesListCmd)
	themesCmd.AddCommand(themesInfoCmd)
	themesCmd.AddCommand(themesInstallCmd)
	themesCmd.AddCommand(themesRemoveCmd)
	themesCmd.AddCommand(themesUpdateCmd)
	rootCmd.AddCommand(themesCmd)
	themesUpdateCmd.Flags().BoolP("check", "c", false, "Check for available updates without installing")
	themesUpdateCmd.Flags().BoolP("all", "a", false, "Update all installed themes")
}
