package betterdiscord

import (
	"log"
	"os"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

func makeDirectory(folder string) error {
	exists := utils.Exists(folder)

	if exists {
		log.Printf("✅ Directory exists: %s", folder)
		return nil
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		log.Printf("❌ Failed to create directory: %s", folder)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Directory created: %s", folder)
	return nil
}

func (i *BDInstall) prepare() error {
	if err := makeDirectory(i.data); err != nil {
		return err
	}
	if err := makeDirectory(i.plugins); err != nil {
		return err
	}
	if err := makeDirectory(i.themes); err != nil {
		return err
	}
	return nil
}

func (i *BDInstall) repair(channel models.DiscordChannel) error {
	channelFolder := filepath.Join(i.data, channel.String())
	pluginsJson := filepath.Join(channelFolder, "plugins.json")

	if !utils.Exists(pluginsJson) {
		log.Printf("✅ No plugins enabled for %s", channel.Name())
		return nil
	}

	if err := os.Remove(pluginsJson); err != nil {
		log.Printf("❌ Unable to remove file %s", pluginsJson)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Plugins disabled for %s", channel.Name())
	return nil
}
