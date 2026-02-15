package betterdiscord

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/betterdiscord/cli/internal/models"
)

type BDInstall struct {
	root          string
	data          string
	asar          string
	plugins       string
	themes        string
	hasDownloaded bool
}

// Root returns the root directory path of the BetterDiscord installation
func (i *BDInstall) Root() string {
	return i.root
}

// Data returns the data directory path
func (i *BDInstall) Data() string {
	return i.data
}

// Asar returns the path to the BetterDiscord asar file
func (i *BDInstall) Asar() string {
	return i.asar
}

// Plugins returns the plugins directory path
func (i *BDInstall) Plugins() string {
	return i.plugins
}

// Themes returns the themes directory path
func (i *BDInstall) Themes() string {
	return i.themes
}

// HasDownloaded returns whether BetterDiscord has been downloaded
func (i *BDInstall) HasDownloaded() bool {
	return i.hasDownloaded
}

// Download downloads the BetterDiscord asar file
func (i *BDInstall) Download() error {
	return i.download()
}

// Prepare creates all necessary directories for BetterDiscord
func (i *BDInstall) Prepare() error {
	return i.prepare()
}

// Repair disables plugins for a specific Discord channel
func (i *BDInstall) Repair(channel models.DiscordChannel) error {
	return i.repair(channel)
}

var lock = &sync.Mutex{}
var globalInstance *BDInstall

func GetInstallation(base ...string) *BDInstall {
	if len(base) == 0 {
		if globalInstance != nil {
			return globalInstance
		}

		lock.Lock()
		defer lock.Unlock()
		if globalInstance != nil {
			return globalInstance
		}

		// Default to user config directory
		configDir, _ := os.UserConfigDir()

		// Handle WSL with Windows home directory
		if os.Getenv("WSL_DISTRO_NAME") != "" && os.Getenv("WIN_HOME") != "" {
			configDir = filepath.Join(os.Getenv("WIN_HOME"), "AppData", "Roaming")
		}

		globalInstance = GetInstallation(configDir)

		return globalInstance
	}

	return New(filepath.Join(base[0], "BetterDiscord"))
}

func New(root string) *BDInstall {
	return &BDInstall{
		root:    root,
		data:    filepath.Join(root, "data"),
		asar:    filepath.Join(root, "data", "betterdiscord.asar"),
		plugins: filepath.Join(root, "plugins"),
		themes:  filepath.Join(root, "themes"),
	}
}
