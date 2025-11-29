package discord

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
var allDiscordInstalls map[models.DiscordChannel][]*DiscordInstall

func GetAllInstalls() map[models.DiscordChannel][]*DiscordInstall {
	var installs = map[models.DiscordChannel][]*DiscordInstall{}

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			installs[result.channel] = append(installs[result.channel], result)
		}
	}

	sortInstalls()

	return installs
}

func GetVersion(proposed string) string {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		if version := versionRegex.FindString(folder); version != "" {
			return version
		}
	}
	return ""
}

func GetChannel(proposed string) models.DiscordChannel {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		for _, channel := range models.Channels {
			if strings.ToLower(folder) == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return models.Stable
}

func GetSuggestedPath(channel models.DiscordChannel) string {
	if len(allDiscordInstalls[channel]) > 0 {
		return allDiscordInstalls[channel][0].corePath
	}
	return ""
}

func AddCustomPath(proposed string) *DiscordInstall {
	result := Validate(proposed)
	if result == nil {
		return nil
	}

	// Check if this already exists in our list and return reference
	index := slices.IndexFunc(allDiscordInstalls[result.channel], func(d *DiscordInstall) bool { return d.corePath == result.corePath })
	if index >= 0 {
		return allDiscordInstalls[result.channel][index]
	}

	allDiscordInstalls[result.channel] = append(allDiscordInstalls[result.channel], result)

	sortInstalls()

	return result
}

func ResolvePath(proposed string) *DiscordInstall {
	for channel := range allDiscordInstalls {
		index := slices.IndexFunc(allDiscordInstalls[channel], func(d *DiscordInstall) bool { return d.corePath == proposed })
		if index >= 0 {
			return allDiscordInstalls[channel][index]
		}
	}

	// If it wasn't found as an existing install, try to add it
	return AddCustomPath(proposed)
}

func sortInstalls() {
	for channel := range allDiscordInstalls {
		slices.SortFunc(allDiscordInstalls[channel], func(a, b *DiscordInstall) int {
			switch {
			case a.version > b.version:
				return -1
			case b.version > a.version:
				return 1
			}
			return 0
		})
	}
}
