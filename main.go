package main

import (
	"log"

	"github.com/betterdiscord/cli/cmd"
)

// Version information set by ldflags during build
var (
	version = "dev"     // Set by: -X main.version=v1.0.0
	commit  = "unknown" // Set by: -X main.commit=abc123
	date    = "unknown" // Set by: -X main.date=2025-02-15...
)

func main() {
	log.SetFlags(0)
	// Initialize version info for cmd package
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
