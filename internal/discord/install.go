package discord

import (
	"path/filepath"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
)

type DiscordInstall struct {
	CorePath  string                `json:"corePath"`
	Channel   models.DiscordChannel `json:"channel"`
	Version   string                `json:"version"`
	IsFlatpak bool                  `json:"isFlatpak"`
	IsSnap    bool                  `json:"isSnap"`
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD() error {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps get their own local BD install
	if discord.IsSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.CorePath, "..", "..", "..", "..")))
	}

	// Make BetterDiscord folders
	output.Println("🛠 Preparing BetterDiscord...")
	if err := bd.Prepare(); err != nil {
		return err
	}
	output.Println("✅ BetterDiscord prepared for install")
	output.Blank()

	// Download and write betterdiscord.asar
	output.Println("📥 Downloading BetterDiscord...")
	if err := bd.Download(); err != nil {
		return err
	}
	output.Println("✅ BetterDiscord downloaded")
	output.Blank()

	// Write injection script to discord_desktop_core/index.js
	output.Println("🔌 Injecting into Discord...")
	if err := discord.inject(bd); err != nil {
		return err
	}
	output.Println("✅ Injection successful")
	output.Blank()

	// Terminate and restart Discord if possible
	output.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	output.Blank()

	return nil
}

// UninstallBD removes BetterDiscord from this Discord installation
func (discord *DiscordInstall) UninstallBD() error {
	output.Println("🧹 Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	output.Blank()

	output.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	output.Blank()

	return nil
}

// RepairBD repairs BetterDiscord for this Discord installation
func (discord *DiscordInstall) RepairBD() error {
	if err := discord.UninstallBD(); err != nil {
		return err
	}

	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps get their own local BD install
	if discord.IsSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.CorePath, "..", "..", "..", "..")))
	}

	if err := bd.Repair(discord.Channel); err != nil {
		return err
	}

	return nil
}
