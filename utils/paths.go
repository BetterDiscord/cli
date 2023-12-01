package utils

import (
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
)

var Roaming string
var BetterDiscord string
var Data string
var Plugins string
var Themes string

func init() {
	var configDir, err = os.UserConfigDir()
	if err != nil {
		return
	}
	Roaming = configDir
	BetterDiscord = path.Join(configDir, "BetterDiscord")
	Data = path.Join(BetterDiscord, "data")
	Plugins = path.Join(BetterDiscord, "plugins")
	Themes = path.Join(BetterDiscord, "themes")
}

func Exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

func DiscordPath(channel string) string {
	var channelName = GetChannelName(channel)

	switch op := runtime.GOOS; op {
	case "windows":
		return ValidatePath(path.Join(os.Getenv("LOCALAPPDATA"), channelName))
	case "darwin":
		return ValidatePath(path.Join("/", "Applications", channelName+".app"))
	case "linux":
		return ValidatePath(path.Join(Roaming, strings.ToLower(channelName)))
	default:
		return ""
	}
}

func ValidatePath(proposed string) string {
	switch op := runtime.GOOS; op {
	case "windows":
		return validateWindows(proposed)
	case "darwin", "linux":
		return validateMacLinux(proposed)
	default:
		return ""
	}
}

func Filter[T any](source []T, filterFunc func(T) bool) (ret []T) {
	var returnArray = []T{}
	for _, s := range source {
		if filterFunc(s) {
			returnArray = append(ret, s)
		}
	}
	return returnArray
}

func validateWindows(proposed string) string {
	var finalPath = ""
	var selected = path.Base(proposed)
	if strings.HasPrefix(selected, "Discord") {

		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return ""
		}

		var candidates = Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && len(strings.Split(file.Name(), ".")) == 3 })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		if len(candidates) == 0 {
			log.Fatal("candidates is zero. do you have the correct branch of discord installed?", ReleaseChannel) //TODO
		}
		var versionDir = candidates[len(candidates)-1].Name()

		// Get core wrap like discord_desktop_core-1
		dFiles, err = os.ReadDir(path.Join(proposed, versionDir, "modules"))
		if err != nil {
			return ""
		}
		candidates = Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		var coreWrap = candidates[len(candidates)-1].Name()

		finalPath = path.Join(proposed, versionDir, "modules", coreWrap, "discord_desktop_core")
	}

	// Use a separate if statement because forcing same-line } else if { is gross
	if strings.HasPrefix(proposed, "app-") {
		var dFiles, err = os.ReadDir(path.Join(proposed, "modules"))
		if err != nil {
			return ""
		}
		var candidates = Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		var coreWrap = candidates[len(candidates)-1].Name()
		finalPath = path.Join(proposed, coreWrap, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// If the path and the asar exist, all good
	if Exists(finalPath) && Exists(path.Join(finalPath, "core.asar")) {
		return finalPath
	}

	return ""
}

func validateMacLinux(proposed string) string {
	if strings.Contains(proposed, "/snap") {
		return ""
	}

	var finalPath = ""
	var selected = path.Base(proposed)
	if strings.HasPrefix(selected, "discord") {
		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return ""
		}

		var candidates = Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && len(strings.Split(file.Name(), ".")) == 3 })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		var versionDir = candidates[len(candidates)-1].Name()
		finalPath = path.Join(proposed, versionDir, "modules", "discord_desktop_core")
	}

	if len(strings.Split(selected, ".")) == 3 {
		finalPath = path.Join(proposed, "modules", "discord_desktop_core")
	}

	if selected == "modules" {
		finalPath = path.Join(proposed, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	if Exists(finalPath) && Exists(path.Join(finalPath, "core.asar")) {
		return finalPath
	}

	return ""
}
