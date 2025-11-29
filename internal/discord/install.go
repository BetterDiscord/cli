package discord

import (
	"log"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
)

type DiscordInstall struct {
	corePath  string                `json:"corePath"`
	channel   models.DiscordChannel `json:"channel"`
	version   string                `json:"version"`
	isFlatpak bool                  `json:"isFlatpak"`
	isSnap    bool                  `json:"isSnap"`
}

func (discord *DiscordInstall) GetPath() string {
	return discord.corePath
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD() error {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps get their own local BD install
	if discord.isSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.corePath, "..", "..", "..", "..")))
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
	log.Printf("## Restarting %s...", discord.channel.Name())
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

	log.Printf("## Restarting %s...", discord.channel.Name())
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
	if discord.isSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.corePath, "..", "..", "..", "..")))
	}

	if err := bd.Repair(discord.channel); err != nil {
		return err
	}

	return nil
}
