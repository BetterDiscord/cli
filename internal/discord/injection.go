package discord

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/output"
)

//go:embed assets/injection.js
var injectionScript string

func (discord *DiscordInstall) inject(bd *betterdiscord.BDInstall) error {
	if discord.IsFlatpak {
		cmd := exec.Command("flatpak", "--user", "override", "com.discordapp."+discord.Channel.Exe(), "--filesystem="+bd.Root())
		if err := cmd.Run(); err != nil {
			output.Printf("❌ Could not give flatpak access to %s\n", bd.Root())
			output.Printf("   %s\n", err.Error())
			return err
		}
	}

	if err := os.WriteFile(filepath.Join(discord.CorePath, "index.js"), []byte(injectionScript), 0755); err != nil {
		output.Printf("❌ Unable to write index.js in %s\n", discord.CorePath)
		output.Printf("   %s\n", err.Error())
		return err
	}

	output.Printf("✅ Injected into %s\n", discord.CorePath)
	return nil
}

func (discord *DiscordInstall) uninject() error {
	indexFile := filepath.Join(discord.CorePath, "index.js")

	contents, err := os.ReadFile(indexFile)

	// First try to check the file, but if there's an issue we try to blindly overwrite below
	if err == nil {
		if !strings.Contains(strings.ToLower(string(contents)), "betterdiscord") {
			output.Printf("✅ No injection found for %s\n", discord.Channel.Name())
			return nil
		}
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0o644); err != nil {
		output.Printf("❌ Unable to write file %s\n", indexFile)
		output.Printf("   %s\n", err.Error())
		return err
	}
	output.Printf("✅ Removed from %s\n", discord.Channel.Name())

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
