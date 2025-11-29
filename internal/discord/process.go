package discord

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/v3/process"
)

func (discord *DiscordInstall) restart() error {
	exeName := discord.getFullExe()

	if running, _ := discord.isRunning(); !running {
		log.Printf("✅ %s not running", discord.channel.Name())
		return nil
	}

	if err := discord.kill(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}

	// Use binary found in killing process
	cmd := exec.Command(exeName)
	if discord.isFlatpak {
		cmd = exec.Command("flatpak", "run", "com.discordapp."+discord.channel.Exe())
	} else if discord.isSnap {
		cmd = exec.Command("snap", "run", discord.channel.Exe())
	}

	// Set working directory to user home
	cmd.Dir, _ = os.UserHomeDir()

	if err := cmd.Start(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Restarted %s", discord.channel.Name())
	return nil
}

func (discord *DiscordInstall) isRunning() (bool, error) {
	name := discord.channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return false, fmt.Errorf("could not list processes")
	}

	// Search for desired processe(s)
	for _, p := range processes {
		n, err := p.Name()

		// Ignore processes requiring Admin/Sudo
		if err != nil {
			continue
		}

		// We found our target return
		if n == name {
			return true, nil
		}
	}

	// If we got here, process was not found
	return false, nil
}

func (discord *DiscordInstall) kill() error {
	name := discord.channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return fmt.Errorf("could not list processes")
	}

	// Search for desired processe(s)
	for _, p := range processes {
		n, err := p.Name()

		// Ignore processes requiring Admin/Sudo
		if err != nil {
			continue
		}

		// We found our target, kill it
		if n == name {
			var killErr = p.Kill()

			// We found it but can't kill it, bail out
			if killErr != nil {
				return killErr
			}
		}
	}

	// If we got here, everything was killed without error
	return nil
}

func (discord *DiscordInstall) getFullExe() string {
	name := discord.channel.Exe()

	var exe = ""
	processes, err := process.Processes()
	if err != nil {
		return exe
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			continue
		}
		if n == name {
			if len(exe) == 0 {
				exe, _ = p.Exe()
			}
		}
	}
	return exe
}
