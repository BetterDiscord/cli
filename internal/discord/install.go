package discord

import (
	"log"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
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
	log.Printf("## Preparing BetterDiscord...")
	if err := bd.Prepare(); err != nil {
		return err
	}
	log.Printf("✅ BetterDiscord prepared for install")
	log.Printf("")

	// Download and write betterdiscord.asar
	log.Printf("## Downloading BetterDiscord...")
	if err := bd.Download(); err != nil {
		return err
	}
	log.Printf("✅ BetterDiscord downloaded")
	log.Printf("")

	// Write injection script to discord_desktop_core/index.js
	log.Printf("## Injecting into Discord...")
	if err := discord.inject(bd); err != nil {
		return err
	}
	log.Printf("✅ Injection successful")
	log.Printf("")

	// Terminate and restart Discord if possible
	log.Printf("## Restarting %s...", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Printf("")

	return nil
}

// UninstallBD removes BetterDiscord from this Discord installation
func (discord *DiscordInstall) UninstallBD() error {
	log.Printf("## Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	log.Printf("")

	log.Printf("## Restarting %s...", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Printf("")

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
