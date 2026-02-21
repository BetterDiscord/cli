package cmd

import (
	"fmt"
	"path"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	uninstallCmd.Flags().StringP("path", "p", "", "Path to a Discord installation")
	uninstallCmd.Flags().StringP("channel", "c", "stable", "Discord release channel (stable|ptb|canary)")
	uninstallCmd.Flags().BoolP("full", "f", false, "Fully uninstall BetterDiscord (uninjects all instances and removes all BetterDiscord folders)")
	uninstallCmd.Flags().BoolP("all", "a", false, "Uninject BetterDiscord from all detected Discord installations")
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls BetterDiscord from your Discord",
	Long:  "Uninstall BetterDiscord by specifying --path/--channel for a single install, --all to uninject from all installs, or --full for complete removal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathFlag, _ := cmd.Flags().GetString("path")
		channelFlag, _ := cmd.Flags().GetString("channel")
		fullFlag, _ := cmd.Flags().GetBool("full")
		allFlag, _ := cmd.Flags().GetBool("all")

		pathProvided := pathFlag != ""
		channelProvided := cmd.Flags().Changed("channel")

		// Validate flag combinations
		if (fullFlag || allFlag) && (pathProvided || channelProvided) {
			return fmt.Errorf("--full and --all cannot be used with --path or --channel")
		}

		if fullFlag && allFlag {
			return fmt.Errorf("--full and --all are mutually exclusive")
		}

		if pathProvided && channelProvided {
			return fmt.Errorf("--path and --channel are mutually exclusive")
		}

		// Full uninstall: all installs, delete all BD folders
		if fullFlag {
			installs := getAllInstalls()

			if err := uninstallAll(installs); err != nil {
				return fmt.Errorf("uninstallation failed: %w", err)
			}

			if err := removeAllBetterDiscord(installs); err != nil {
				return fmt.Errorf("folder removal failed: %w", err)
			}

			output.Println("✅ BetterDiscord fully uninstalled from all Discord instances")
			return nil
		}

		// Uninject all: all installs, no deletion
		if allFlag {
			installs := getAllInstalls()

			if err := uninstallAll(installs); err != nil {
				return fmt.Errorf("uninstallation failed: %w", err)
			}

			output.Println("✅ BetterDiscord uninjected from all Discord instances")
			return nil
		}

		// Default: uninject single install
		var install *discord.DiscordInstall

		if pathProvided {
			install = discord.ResolvePath(pathFlag)
			if install == nil {
				return fmt.Errorf("could not find a valid Discord installation at %s", pathFlag)
			}
		} else {
			channel := models.ParseChannel(channelFlag)
			resolvedPath := discord.GetSuggestedPath(channel)
			install = discord.ResolvePath(resolvedPath)
			if install == nil {
				return fmt.Errorf("could not find a valid %s installation to uninstall", channelFlag)
			}
		}

		if err := install.UninstallBD(); err != nil {
			return fmt.Errorf("uninstallation failed: %w", err)
		}

		output.Printf("✅ BetterDiscord uninstalled from %s\n", path.Dir(install.CorePath))
		return nil
	},
}

// getAllInstalls returns all detected Discord installations without duplicates.
func getAllInstalls() []*discord.DiscordInstall {
	installsMap := discord.GetAllInstalls()
	seen := map[string]bool{}
	var installs []*discord.DiscordInstall

	// Flatten the map of installs and filter out duplicates based on CorePath
	// Honestly, probably should have just returned a flat list from GetAllInstalls in the first place, but whatever
	// And also the chance of actually having duplicates is pretty much zero, but this is just in case
	// If you are reading this and you do have duplicates, please tell me because that would be very interesting and I would like to know how that happened
	// If you are reading this and confused by these comments, hi, I'm Zerebos, the author of this code, and I just wanted to have a little fun here, hope you don't mind
	for _, list := range installsMap {
		for _, inst := range list {
			if inst == nil {
				continue
			}
			if seen[inst.CorePath] {
				continue
			}
			seen[inst.CorePath] = true
			installs = append(installs, inst)
		}
	}

	return installs
}

func uninstallAll(installs []*discord.DiscordInstall) error {
	var firstErr error
	for _, inst := range installs {
		if err := inst.UninstallBD(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			output.Printf("❌ Failed to uninstall from %s\n", path.Dir(inst.CorePath))
			output.Printf("   %s\n", err.Error())
		}
	}
	return firstErr
}

func removeAllBetterDiscord(installs []*discord.DiscordInstall) error {
	roots := map[string]*betterdiscord.BDInstall{}

	// This is actually a case where duplicates will happen, because
	// multiple Discord installs can share the same BetterDiscord root,
	// so we need to filter them out to avoid trying to delete the same
	// folder multiple times
	for _, inst := range installs {
		bd := inst.GetBetterDiscordInstall()
		if bd == nil {
			continue
		}
		roots[bd.Root()] = bd
	}

	var firstErr error
	for _, bd := range roots {
		if err := bd.RemoveAll(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}
