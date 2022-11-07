package utils

import (
	"os"
	"path"
	"runtime"
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
	case "darwin":
		return validateMac(proposed)
	case "linux":
		return validateLinux(proposed)
	default:
		return ""
	}
}

func validateWindows(proposed string) string {
	var finalPath = ""
	var selected = path.Base(proposed)
	if strings.HasPrefix(selected, "Discord") {
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return ""
		}
		var appDir = ""
		for _, file := range dFiles {
			if !file.IsDir() || !strings.HasPrefix(file.Name(), "app-") {
				continue
			}
			if file.Name() > appDir {
				appDir = file.Name()
			}
		}
		finalPath = path.Join(proposed, appDir, "resources")
	}

	// Use a separate if statement because forcing same-line } else if { is gross
	if strings.HasPrefix(proposed, "app-") {
		finalPath = path.Join(proposed, "resources")
	}

	if selected == "resources" {
		finalPath = proposed
	}

	if Exists(finalPath) {
		return finalPath
	}

	return ""
}

func validateMac(proposed string) string {
	var finalPath = ""
	var selected = path.Base(proposed)
	if strings.HasPrefix(selected, "Discord") && strings.HasSuffix(selected, ".app") {
		finalPath = path.Join(proposed, "Contents", "Resources")
	}

	if selected == "Contents" {
		finalPath = path.Join(proposed, "Resources")
	}

	if selected == "Resources" {
		finalPath = proposed
	}

	if Exists(finalPath) {
		return finalPath
	}

	return ""
}

func validateLinux(proposed string) string {
	if strings.Contains(proposed, "/snap") {
		return ""
	}

	var finalPath = ""
	var selected = path.Base(proposed)
	if strings.HasPrefix(selected, "discord") {
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return ""
		}
		var versionDir = ""
		for _, file := range dFiles {
			if split := strings.Split(file.Name(), "."); !file.IsDir() || len(split) != 3 {
				continue
			}
			if file.Name() > versionDir {
				versionDir = file.Name()
			}
		}
		finalPath = path.Join(proposed, versionDir, "modules", "discord_desktop_core")
	}

	if split := strings.Split(selected, "."); len(split) == 3 {
		finalPath = path.Join(proposed, "modules", "discord_desktop_core")
	}

	if selected == "modules" {
		finalPath = path.Join(proposed, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	if Exists(finalPath) {
		return finalPath
	}

	return ""
}
