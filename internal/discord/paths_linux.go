package discord

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

		// WSL. Data is stored under the Windows user's AppData folder.
		// Example: `/mnt/c/Users/Username/AppData/Local/DiscordCanary`.
		// Core: `/mnt/c/Users/Username/AppData/Local/DiscordCanary/app-1.0.9218/modules/discord_desktop_core-1/discord_desktop_core core.asar`.
		filepath.Join(os.Getenv("WIN_HOME"), "AppData", "Local", "{CHANNEL}"),
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

/**
 * Currently nearly the same as darwin validation however
 * it is kept separate in case of future changes to
 * either system, it is likely that linux will require
 * more advanced validation for snap and flatpak.
 */
func Validate(proposed string) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)
	if strings.HasPrefix(strings.ToLower(selected), "discord") {
		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return nil
		}

		var candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && versionRegex.MatchString(file.Name()) })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		var versionDir = candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core")

		// WSL installs have an extra layer
		if (os.Getenv("WSL_DISTRO_NAME") != "" && os.Getenv("WIN_HOME") != "") && strings.Contains(proposed, "AppData") {
			finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")
		}
	}

	if len(strings.Split(selected, ".")) == 3 {
		finalPath = filepath.Join(proposed, "modules", "discord_desktop_core")
	}

	if selected == "modules" {
		finalPath = filepath.Join(proposed, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// If the path and the asar exist, all good
	if utils.Exists(finalPath) && utils.Exists(filepath.Join(finalPath, "core.asar")) {
		return &DiscordInstall{
			CorePath:  finalPath,
			Channel:   GetChannel(finalPath),
			Version:   GetVersion(finalPath),
			IsFlatpak: strings.Contains(finalPath, "com.discordapp."),
			IsSnap:    strings.Contains(finalPath, "snap/"),
		}
	}

	return nil
}
