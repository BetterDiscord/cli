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
	paths := []string{
		filepath.Join(config, "{channel}"),
	}

	for _, channel := range models.Channels {
		for _, path := range paths {
			folder := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "")
			searchPaths = append(
				searchPaths,
				strings.ReplaceAll(path, "{channel}", folder),
			)
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

/**
 * Currently nearly the same as linux validation however
 * it is kept separate in case of future changes to
 * either system, it is likely that linux will require
 * more advanced validation for snap and flatpak.
 */
func Validate(proposed string) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)
	if strings.HasPrefix(selected, "discord") {
		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return nil
		}

		var candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && versionRegex.MatchString(file.Name()) })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		var versionDir = candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core")
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
			IsFlatpak: false,
			IsSnap:    false,
		}
	}

	return nil
}
