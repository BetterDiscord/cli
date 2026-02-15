package discord

import (
	"path/filepath"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Version in middle of path",
			path:     "/usr/share/discord/0.0.35/modules",
			expected: "0.0.35",
		},
		{
			name:     "Version at end of path",
			path:     "/home/user/.config/discord/0.0.36",
			expected: "0.0.36",
		},
		{
			name:     "Multiple versions (should return first)",
			path:     "/usr/share/1.2.3/discord/0.0.35/modules",
			expected: "1.2.3",
		},
		{
			name:     "No version in path",
			path:     "/usr/share/discord/modules",
			expected: "",
		},
		{
			name:     "Windows-style path with version",
			path:     "C:\\Users\\User\\AppData\\Local\\Discord\\app-1.0.9012",
			expected: "1.0.9012",
		},
		{
			name:     "Version with many digits",
			path:     "/opt/discord/123.456.789/core",
			expected: "123.456.789",
		},
		{
			name:     "Empty path",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVersion(tt.path)
			if result != tt.expected {
				t.Errorf("GetVersion(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetChannel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected models.DiscordChannel
	}{
		{
			name:     "Stable in path (lowercase)",
			path:     "/usr/share/discord/modules",
			expected: models.Stable,
		},
		{
			name:     "Canary in path (lowercase)",
			path:     "/usr/share/discordcanary/modules",
			expected: models.Canary,
		},
		{
			name:     "PTB in path (lowercase)",
			path:     "/usr/share/discordptb/modules",
			expected: models.PTB,
		},

		{
			name:     "DiscordCanary without space",
			path:     "/home/user/.config/DiscordCanary/modules",
			expected: models.Canary,
		},
		{
			name:     "DiscordPTB without space",
			path:     "/home/user/.config/DiscordPTB/modules",
			expected: models.PTB,
		},
		{
			name:     "No channel identifier defaults to Stable",
			path:     "/some/random/path/modules",
			expected: models.Stable,
		},
		{
			name:     "Multiple Discord mentions (first wins)",
			path:     filepath.Join("discordcanary", "discord", "modules"),
			expected: models.Canary,
		},
		{
			name:     "Empty path defaults to Stable",
			path:     "",
			expected: models.Stable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetChannel(tt.path)
			if result != tt.expected {
				t.Errorf("GetChannel(%s) = %v (%s), expected %v (%s)",
					tt.path, result, result.String(), tt.expected, tt.expected.String())
			}
		})
	}
}

func TestGetChannel_CaseInsensitive(t *testing.T) {
	tests := []struct {
		path     string
		expected models.DiscordChannel
	}{
		{"/usr/share/DISCORD/modules", models.Stable},
		{"/usr/share/Discord/modules", models.Stable},
		{"/usr/share/DISCORDCANARY/modules", models.Canary},
		{"/usr/share/DiscordCanary/modules", models.Canary},
		{"/usr/share/DISCORDPTB/modules", models.PTB},
		{"/usr/share/DiscordPTB/modules", models.PTB},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetChannel(tt.path)
			if result != tt.expected {
				t.Errorf("GetChannel(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetSuggestedPath(t *testing.T) {
	// Reset allDiscordInstalls for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Test empty installs
	result := GetSuggestedPath(models.Stable)
	if result != "" {
		t.Errorf("GetSuggestedPath with no installs should return empty string, got %s", result)
	}

	// Add some test installs
	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{CorePath: "/usr/share/discord/0.0.35", Version: "0.0.35"},
		{CorePath: "/usr/share/discord/0.0.34", Version: "0.0.34"},
	}

	allDiscordInstalls[models.Canary] = []*DiscordInstall{
		{CorePath: "/usr/share/discord-canary/0.0.200", Version: "0.0.200"},
	}

	// Test that it returns the first install
	stableResult := GetSuggestedPath(models.Stable)
	if stableResult != "/usr/share/discord/0.0.35" {
		t.Errorf("GetSuggestedPath(Stable) = %s, expected /usr/share/discord/0.0.35", stableResult)
	}

	canaryResult := GetSuggestedPath(models.Canary)
	if canaryResult != "/usr/share/discord-canary/0.0.200" {
		t.Errorf("GetSuggestedPath(Canary) = %s, expected /usr/share/discord-canary/0.0.200", canaryResult)
	}

	// Test channel with no installs
	ptbResult := GetSuggestedPath(models.PTB)
	if ptbResult != "" {
		t.Errorf("GetSuggestedPath(PTB) with no PTB installs should return empty string, got %s", ptbResult)
	}
}

func TestAddCustomPath(t *testing.T) {
	// This test is limited because Validate() depends on OS-specific paths
	// We're mainly testing the logic around adding and deduplication

	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Test with invalid path (will return nil since Validate will fail)
	result := AddCustomPath("/nonexistent/invalid/path")
	if result != nil {
		t.Error("AddCustomPath with invalid path should return nil")
	}

	// Further testing would require mocking the Validate function
	// or setting up actual Discord installation directories
}

func TestResolvePath(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add a test install
	testInstall := &DiscordInstall{
		CorePath: "/test/discord/path",
		Channel:  models.Stable,
		Version:  "1.0.0",
	}
	allDiscordInstalls[models.Stable] = []*DiscordInstall{testInstall}

	// Test resolving existing path
	result := ResolvePath("/test/discord/path")
	if result != testInstall {
		t.Error("ResolvePath should return the existing install")
	}

	// Test resolving non-existent path (will try AddCustomPath and likely return nil)
	result2 := ResolvePath("/nonexistent/path")
	if result2 != nil {
		// This might succeed or fail depending on whether Validate passes
		// In most test environments, it should return nil
		t.Log("ResolvePath returned non-nil for non-existent path (may be valid in some environments)")
	}
}

func TestSortInstalls(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add unsorted installs
	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{CorePath: "/path1", Version: "0.0.34", Channel: models.Stable},
		{CorePath: "/path2", Version: "0.0.36", Channel: models.Stable},
		{CorePath: "/path3", Version: "0.0.35", Channel: models.Stable},
	}

	// Sort them
	sortInstalls()

	// Verify sorted in descending order by version
	installs := allDiscordInstalls[models.Stable]
	if len(installs) != 3 {
		t.Fatalf("Expected 3 installs, got %d", len(installs))
	}

	if installs[0].Version != "0.0.36" {
		t.Errorf("First install should have version 0.0.36, got %s", installs[0].Version)
	}
	if installs[1].Version != "0.0.35" {
		t.Errorf("Second install should have version 0.0.35, got %s", installs[1].Version)
	}
	if installs[2].Version != "0.0.34" {
		t.Errorf("Third install should have version 0.0.34, got %s", installs[2].Version)
	}
}

func TestSortInstalls_MultipleChannels(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add unsorted installs for multiple channels
	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{CorePath: "/stable1", Version: "1.0.0", Channel: models.Stable},
		{CorePath: "/stable2", Version: "1.0.2", Channel: models.Stable},
	}

	allDiscordInstalls[models.Canary] = []*DiscordInstall{
		{CorePath: "/canary1", Version: "0.0.100", Channel: models.Canary},
		{CorePath: "/canary2", Version: "0.0.150", Channel: models.Canary},
		{CorePath: "/canary3", Version: "0.0.125", Channel: models.Canary},
	}

	// Sort them
	sortInstalls()

	// Verify Stable channel is sorted
	stableInstalls := allDiscordInstalls[models.Stable]
	if stableInstalls[0].Version != "1.0.2" {
		t.Errorf("Stable: First version should be 1.0.2, got %s", stableInstalls[0].Version)
	}
	if stableInstalls[1].Version != "1.0.0" {
		t.Errorf("Stable: Second version should be 1.0.0, got %s", stableInstalls[1].Version)
	}

	// Verify Canary channel is sorted
	canaryInstalls := allDiscordInstalls[models.Canary]
	if canaryInstalls[0].Version != "0.0.150" {
		t.Errorf("Canary: First version should be 0.0.150, got %s", canaryInstalls[0].Version)
	}
	if canaryInstalls[1].Version != "0.0.125" {
		t.Errorf("Canary: Second version should be 0.0.125, got %s", canaryInstalls[1].Version)
	}
	if canaryInstalls[2].Version != "0.0.100" {
		t.Errorf("Canary: Third version should be 0.0.100, got %s", canaryInstalls[2].Version)
	}
}
