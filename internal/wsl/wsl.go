package wsl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var once sync.Once
var info *WSLInfo

type WSLInfo struct {
	IsWSL         bool   // True if running under WSL1 or WSL2
	DistroName    string // Value of WSL_DISTRO_NAME
	KernelVersion string // Contents of /proc/version
	InteropPath   string // Value of WSL_INTEROP
	WindowsHome   string // Windows home directory as a WSL path
}

// loadWSLInfo detects WSL environment details and caches them in the info variable.
// This function is safe to call multiple times but will only execute once.
// It checks for WSL using both environment variables and kernel signatures, and attempts
// to determine the Windows home directory if running under WSL.
//
// This is called lazily by public helper functions to avoid unnecessary overhead
// of init() in non-WSL environments.
func loadWSLInfo() {
	once.Do(func() {
		i := &WSLInfo{}

		// Detect WSL via environment variable
		if dn := os.Getenv("WSL_DISTRO_NAME"); dn != "" {
			i.IsWSL = true
			i.DistroName = dn
		}

		// Kernel signature fallback
		if data, err := os.ReadFile("/proc/version"); err == nil {
			ver := strings.ToLower(string(data))
			i.KernelVersion = ver

			// Check for "microsoft" in kernel version to detect WSL if env var is missing
			if strings.Contains(ver, "microsoft") {
				i.IsWSL = true
			}
		}

		// Interop path (could be useful someday)
		i.InteropPath = os.Getenv("WSL_INTEROP")

		// Windows home directory (only if WSL)
		if i.IsWSL {
			home, err := getWindowsHomePath()
			if err == nil {
				i.WindowsHome = home
			}
		}

		info = i
	})
}

// Info returns cached WSL environment information.
func Info() *WSLInfo {
	loadWSLInfo()
	return info
}

// IsWSL returns true if running under WSL.
func IsWSL() bool {
	loadWSLInfo()
	return info.IsWSL
}

// WindowsHome returns the Windows user's home directory as a WSL path.
func WindowsHome() (string, error) {
	loadWSLInfo()
	if info.WindowsHome == "" {
		return "", fmt.Errorf("unable to determine Windows home directory")
	}
	return info.WindowsHome, nil
}

// ToWSLPath converts a Windows path (C:\Users\Me) → /mnt/c/Users/Me.
func ToWSLPath(winPath string) (string, error) {
	out, err := exec.Command("wslpath", "-u", winPath).Output()
	if err != nil {
		return "", fmt.Errorf("wslpath failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ToWindowsPath converts a WSL path (/mnt/c/Users/Me) → C:\Users\Me.
func ToWindowsPath(wslPath string) (string, error) {
	out, err := exec.Command("wslpath", "-w", wslPath).Output()
	if err != nil {
		return "", fmt.Errorf("wslpath failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ExecWindows runs a Windows command and returns stdout.
func ExecWindows(command string) (string, error) {
	// Use cmd.exe to run the command
	cmd := exec.Command("cmd.exe", "/c", command)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.exe failed: %w", err)
	}

	// Clean CRLF from Windows output
	return strings.TrimSpace(strings.ReplaceAll(string(out), "\r", "")), nil
}

func getWindowsHomePath() (string, error) {
	// Use ExecWindows helper
	winPath, err := ExecWindows("echo %USERPROFILE%")
	if err != nil {
		return "", err
	}

	// Convert to WSL path
	return ToWSLPath(winPath)
}
