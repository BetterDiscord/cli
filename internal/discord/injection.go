package discord

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/betterdiscord"
)

//go:embed assets/injection.js
var injectionScript string

func (discord *DiscordInstall) inject(bd *betterdiscord.BDInstall) error {
	if discord.IsFlatpak {
		cmd := exec.Command("flatpak", "--user", "override", "com.discordapp."+discord.Channel.Exe(), "--filesystem="+bd.Root())
		if err := cmd.Run(); err != nil {
			log.Printf("❌ Could not give flatpak access to %s", bd.Root())
			log.Printf("❌ %s", err.Error())
			return err
		}
	}

	if err := os.WriteFile(filepath.Join(discord.CorePath, "index.js"), []byte(injectionScript), 0755); err != nil {
		log.Printf("❌ Unable to write index.js in %s", discord.CorePath)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Injected into %s", discord.CorePath)
	return nil
}

func (discord *DiscordInstall) uninject() error {
	indexFile := filepath.Join(discord.CorePath, "index.js")

	contents, err := os.ReadFile(indexFile)

	// First try to check the file, but if there's an issue we try to blindly overwrite below
	if err == nil {
		if !strings.Contains(strings.ToLower(string(contents)), "betterdiscord") {
			log.Printf("✅ No injection found for %s", discord.Channel.Name())
			return nil
		}
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0o644); err != nil {
		log.Printf("❌ Unable to write file %s", indexFile)
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Removed from %s", discord.Channel.Name())

	return nil
}

// TODO: consider putting this in the betterdiscord package
func (discord *DiscordInstall) IsInjected() bool {
	indexFile := filepath.Join(discord.CorePath, "index.js")
	contents, err := os.ReadFile(indexFile)
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(contents))
	return strings.Contains(lower, "betterdiscord")
}
