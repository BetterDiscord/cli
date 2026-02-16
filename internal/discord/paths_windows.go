package discord

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
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
	return validateWindowsStyleInstall(proposed)
}
