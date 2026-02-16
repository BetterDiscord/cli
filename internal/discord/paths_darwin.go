package discord

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
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

func Validate(proposed string) *DiscordInstall {
	return validateUnixStyleInstall(proposed, false, false)
}
