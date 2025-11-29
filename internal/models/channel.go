package models

import (
	"runtime"
	"strings"
)

// DiscordChannel represents a Discord release channel (Stable, PTB, Canary)
type DiscordChannel int

const (
	Stable DiscordChannel = iota
	Canary
	PTB
)

// All available Discord channels
var Channels = []DiscordChannel{Stable, Canary, PTB}

// Used for logging, etc
func (channel DiscordChannel) String() string {
	switch channel {
	case Stable:
		return "stable"
	case Canary:
		return "canary"
	case PTB:
		return "ptb"
	}
	return ""
}

// Used for user display
func (channel DiscordChannel) Name() string {
	switch channel {
	case Stable:
		return "Discord"
	case Canary:
		return "Discord Canary"
	case PTB:
		return "Discord PTB"
	}
	return ""
}

// Exe returns the executable name for the release channel
func (channel DiscordChannel) Exe() string {
	name := channel.Name()

	if runtime.GOOS != "darwin" {
		name = strings.ReplaceAll(name, " ", "")
	}

	if runtime.GOOS == "windows" {
		name = name + ".exe"
	}

	return name
}

// ParseChannel converts a string input to a DiscordChannel type
func ParseChannel(input string) DiscordChannel {
	switch strings.ToLower(input) {
	case "stable":
		return Stable
	case "canary":
		return Canary
	case "ptb":
		return PTB
	}
	return Stable
}

// Used by Wails for type serialization
func (channel DiscordChannel) TSName() string {
	return strings.ToUpper(channel.String())
}
