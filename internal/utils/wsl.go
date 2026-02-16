package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	data, err := os.ReadFile("/proc/version")
	if err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft") {
		return true
	}

	return false
}

func WindowsHome() (string, error) {
	// Run Windows command
	cmd := exec.Command("cmd.exe", "/c", "echo %USERPROFILE%")
	raw, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.exe failed: %w", err)
	}

	// Clean CRLF
	winPath := strings.TrimSpace(strings.ReplaceAll(string(raw), "\r", ""))

	// Convert to WSL path
	cmd2 := exec.Command("wslpath", "-u", winPath)
	out, err := cmd2.Output()
	if err != nil {
		return "", fmt.Errorf("wslpath failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
