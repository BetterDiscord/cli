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
	paths := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "{channel}"),
		filepath.Join(os.Getenv("PROGRAMDATA"), os.Getenv("USERNAME"), "{channel}"),
	}

	for _, channel := range models.Channels {
		for _, path := range paths {
			searchPaths = append(
				searchPaths,
				strings.ReplaceAll(path, "{channel}", strings.ReplaceAll(channel.Name(), " ", "")),
			)
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

func Validate(proposed string) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)

	if strings.HasPrefix(selected, "Discord") {

		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return nil
		}

		var candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && len(strings.Split(file.Name(), ".")) == 3 })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		var versionDir = candidates[len(candidates)-1].Name()

		// Get core wrap like discord_desktop_core-1
		dFiles, err = os.ReadDir(filepath.Join(proposed, versionDir, "modules"))
		if err != nil {
			return nil
		}
		candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		var coreWrap = candidates[len(candidates)-1].Name()

		finalPath = filepath.Join(proposed, versionDir, "modules", coreWrap, "discord_desktop_core")
	}

	// Use a separate if statement because forcing same-line } else if { is gross
	if strings.HasPrefix(selected, "app-") {
		var dFiles, err = os.ReadDir(filepath.Join(proposed, "modules"))
		if err != nil {
			return nil
		}

		var candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		var coreWrap = candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, "modules", coreWrap, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// If the path and the asar exist, all good
	if utils.Exists(finalPath) && utils.Exists(filepath.Join(finalPath, "core.asar")) {
		return &DiscordInstall{
			corePath:  finalPath,
			channel:   GetChannel(finalPath),
			version:   GetVersion(finalPath),
			isFlatpak: false,
			isSnap:    false,
		}
	}

	return nil
}
