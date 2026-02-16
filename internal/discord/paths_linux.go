package discord

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

func init() {
	config, _ := os.UserConfigDir()
	home, _ := os.UserHomeDir()
	paths := []string{
		// Native. Data is stored under `~/.config`.
		// Example: `~/.config/discordcanary`.
		// Core: `~/.config/discordcanary/0.0.90/modules/discord_desktop_core/core.asar`.
		filepath.Join(config, "{channel}"),

		// Flatpak. These user data paths are universal for all Flatpak installations on all machines.
		// Example: `.var/app/com.discordapp.DiscordCanary/config/discordcanary`.
		// Core: `.var/app/com.discordapp.DiscordCanary/config/discordcanary/0.0.90/modules/discord_desktop_core/core.asar`
		filepath.Join(home, ".var", "app", "com.discordapp.{CHANNEL}", "config", "{channel}"),

		// Snap. Just like with Flatpaks, these paths are universal for all Snap installations.
		// Example: `snap/discord/current/.config/discord`.
		// Example: `snap/discord-canary/current/.config/discordcanary`.
		// Core: `snap/discord-canary/current/.config/discordcanary/0.0.90/modules/discord_desktop_core/core.asar`.
		// NOTE: Snap user data always exists, even when the Snap isn't mounted/running.
		filepath.Join(home, "snap", "{channel-}", "current", ".config", "{channel}"),
	}

	if utils.IsWSL() {
		winHome, err := utils.WindowsHome()
		if err == nil && winHome != "" {
			// WSL. Data is stored under the Windows user's AppData folder.
			// Example: `/mnt/c/Users/Username/AppData/Local/DiscordCanary`.
			// Core: `/mnt/c/Users/Username/AppData/Local/DiscordCanary/app-1.0.9218/modules/discord_desktop_core-1/discord_desktop_core core.asar`.
			paths = append(paths, filepath.Join(winHome, "AppData", "Local", "{CHANNEL}"))
		}
	}

	for _, channel := range models.Channels {
		for _, path := range paths {
			upper := strings.ReplaceAll(channel.Name(), " ", "")
			lower := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "")
			dash := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "-")
			folder := strings.ReplaceAll(path, "{CHANNEL}", upper)
			folder = strings.ReplaceAll(folder, "{channel}", lower)
			folder = strings.ReplaceAll(folder, "{channel-}", dash)
			searchPaths = append(searchPaths, folder)
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

// Validate validates a Discord installation path on Linux.
// For WSL environments, it uses Windows-style validation.
// For native Linux, it detects Flatpak and Snap installations.
func Validate(proposed string) *DiscordInstall {
	if utils.IsWSL() {
		return validateWindowsStyleInstall(proposed)
	}

	// Native Linux validation with Flatpak and Snap detection
	return validateUnixStyleInstall(proposed, true, true)
}
