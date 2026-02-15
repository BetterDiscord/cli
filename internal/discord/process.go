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
		log.Printf("✅ %s not running", discord.Channel.Name())
		return nil
	}

	if err := discord.kill(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.Channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}

	// Determine command based on installation type
	var cmd *exec.Cmd
	if discord.IsFlatpak {
		cmd = exec.Command("flatpak", "run", "com.discordapp."+discord.Channel.Exe())
	} else if discord.IsSnap {
		cmd = exec.Command("snap", "run", discord.Channel.Exe())
	} else {
		// Use binary found in killing process for non-Flatpak/Snap installs
		if exeName == "" {
			log.Printf("❌ Unable to restart %s, please do so manually!", discord.Channel.Name())
			return fmt.Errorf("could not determine executable path for %s", discord.Channel.Name())
		}
		cmd = exec.Command(exeName)
	}

	// Set working directory to user home
	cmd.Dir, _ = os.UserHomeDir()

	if err := cmd.Start(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.Channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Restarted %s", discord.Channel.Name())
	return nil
}

func (discord *DiscordInstall) isRunning() (bool, error) {
	name := discord.Channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return false, fmt.Errorf("could not list processes")
	}

	// Search for desired process(es)
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
	name := discord.Channel.Exe()
	processes, err := process.Processes()

	// If we can't even list processes, bail out
	if err != nil {
		return fmt.Errorf("could not list processes")
	}

	// Search for desired process(es)
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
	name := discord.Channel.Exe()

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
