package models

import (
	"runtime"
	"testing"
)

func TestDiscordChannel_String(t *testing.T) {
	tests := []struct {
		channel  DiscordChannel
		expected string
	}{
		{Stable, "stable"},
		{Canary, "canary"},
		{PTB, "ptb"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.channel.String()
			if result != tt.expected {
				t.Errorf("String() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestDiscordChannel_Name(t *testing.T) {
	tests := []struct {
		channel  DiscordChannel
		expected string
	}{
		{Stable, "Discord"},
		{Canary, "Discord Canary"},
		{PTB, "Discord PTB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.channel.Name()
			if result != tt.expected {
				t.Errorf("Name() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestDiscordChannel_Exe(t *testing.T) {
	tests := []struct {
		name     string
		channel  DiscordChannel
		goos     string
		expected string
	}{
		{
			name:     "Stable on Linux",
			channel:  Stable,
			goos:     "linux",
			expected: "Discord",
		},
		{
			name:     "Canary on Linux",
			channel:  Canary,
			goos:     "linux",
			expected: "DiscordCanary",
		},
		{
			name:     "PTB on Linux",
			channel:  PTB,
			goos:     "linux",
			expected: "DiscordPTB",
		},
		{
			name:     "Stable on Darwin",
			channel:  Stable,
			goos:     "darwin",
			expected: "Discord",
		},
		{
			name:     "Canary on Darwin",
			channel:  Canary,
			goos:     "darwin",
			expected: "Discord Canary",
		},
		{
			name:     "PTB on Darwin",
			channel:  PTB,
			goos:     "darwin",
			expected: "Discord PTB",
		},
		{
			name:     "Stable on Windows",
			channel:  Stable,
			goos:     "windows",
			expected: "Discord.exe",
		},
		{
			name:     "Canary on Windows",
			channel:  Canary,
			goos:     "windows",
			expected: "DiscordCanary.exe",
		},
		{
			name:     "PTB on Windows",
			channel:  PTB,
			goos:     "windows",
			expected: "DiscordPTB.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test will check the actual runtime.GOOS
			// In a real test environment, you might want to mock this
			// For now, we'll just test the actual platform
			if runtime.GOOS != tt.goos {
				t.Skipf("Skipping test for %s on %s", tt.goos, runtime.GOOS)
			}

			result := tt.channel.Exe()
			if result != tt.expected {
				t.Errorf("Exe() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestDiscordChannel_Exe_CurrentPlatform(t *testing.T) {
	// Test that Exe() returns something reasonable for the current platform
	channels := []DiscordChannel{Stable, Canary, PTB}

	for _, channel := range channels {
		result := channel.Exe()

		// Basic validation
		if result == "" {
			t.Errorf("Exe() returned empty string for %s", channel.Name())
		}

		// Platform-specific checks
		switch runtime.GOOS {
		case "windows":
			if result[len(result)-4:] != ".exe" {
				t.Errorf("Exe() on Windows should end with .exe, got %s", result)
			}
		case "darwin":
			// On macOS, names should contain spaces for Canary and PTB
			if channel == Canary && result != "Discord Canary" {
				t.Errorf("Exe() on macOS for Canary = %s, expected 'Discord Canary'", result)
			}
			if channel == PTB && result != "Discord PTB" {
				t.Errorf("Exe() on macOS for PTB = %s, expected 'Discord PTB'", result)
			}
		default:
			// On Linux, names should not contain spaces
			if channel == Canary && result != "DiscordCanary" {
				t.Errorf("Exe() on Linux for Canary = %s, expected 'DiscordCanary'", result)
			}
			if channel == PTB && result != "DiscordPTB" {
				t.Errorf("Exe() on Linux for PTB = %s, expected 'DiscordPTB'", result)
			}
		}
	}
}

func TestParseChannel(t *testing.T) {
	tests := []struct {
		input    string
		expected DiscordChannel
	}{
		{"stable", Stable},
		{"Stable", Stable},
		{"STABLE", Stable},
		{"canary", Canary},
		{"Canary", Canary},
		{"CANARY", Canary},
		{"ptb", PTB},
		{"PTB", PTB},
		{"Ptb", PTB},
		{"invalid", Stable}, // Default to Stable
		{"", Stable},        // Default to Stable
		{"discord", Stable}, // Unknown input defaults to Stable
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseChannel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseChannel(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDiscordChannel_TSName(t *testing.T) {
	tests := []struct {
		channel  DiscordChannel
		expected string
	}{
		{Stable, "STABLE"},
		{Canary, "CANARY"},
		{PTB, "PTB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.channel.TSName()
			if result != tt.expected {
				t.Errorf("TSName() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestChannelsConstant(t *testing.T) {
	// Verify that Channels contains all expected channels
	if len(Channels) != 3 {
		t.Errorf("Channels should contain 3 channels, got %d", len(Channels))
	}

	expectedChannels := []DiscordChannel{Stable, Canary, PTB}
	for i, expected := range expectedChannels {
		if Channels[i] != expected {
			t.Errorf("Channels[%d] = %v, expected %v", i, Channels[i], expected)
		}
	}
}

func TestDiscordChannel_EnumValues(t *testing.T) {
	// Test that the enum values are distinct
	if Stable == Canary {
		t.Error("Stable and Canary should have different values")
	}
	if Stable == PTB {
		t.Error("Stable and PTB should have different values")
	}
	if Canary == PTB {
		t.Error("Canary and PTB should have different values")
	}

	// Test that values are sequential starting from 0
	if Stable != 0 {
		t.Errorf("Stable should be 0, got %d", Stable)
	}
	if Canary != 1 {
		t.Errorf("Canary should be 1, got %d", Canary)
	}
	if PTB != 2 {
		t.Errorf("PTB should be 2, got %d", PTB)
	}
}
