package betterdiscord

import (
	"os"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

func makeDirectory(folder string) error {
	exists := utils.Exists(folder)

	if exists {
		output.Printf("✅ Directory exists: %s\n", folder)
		return nil
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		output.Printf("❌ Failed to create directory: %s\n", folder)
		output.Printf("   %s\n", err.Error())
		return err
	}

	output.Printf("✅ Directory created: %s\n", folder)
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
		output.Printf("✅ No plugins enabled for %s\n", channel.Name())
		return nil
	}

	if err := os.Remove(pluginsJson); err != nil {
		output.Printf("❌ Unable to remove file %s\n", pluginsJson)
		output.Printf("   %s\n", err.Error())
		return err
	}

	output.Printf("✅ Plugins disabled for %s\n", channel.Name())
	return nil
}
