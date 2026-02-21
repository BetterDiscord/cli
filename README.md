# BetterDiscord CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/BetterDiscord/cli)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/BetterDiscord/cli)](https://github.com/BetterDiscord/cli/releases)
[![License](https://img.shields.io/github/license/BetterDiscord/cli)](LICENSE)
[![npm](https://img.shields.io/npm/v/@betterdiscord/cli)](https://www.npmjs.com/package/@betterdiscord/cli)

A cross-platform command-line interface for installing, updating, and managing [BetterDiscord](https://betterdiscord.app/).

## Features

- 🚀 Easy installation and uninstallation of BetterDiscord
- 🔄 Support for multiple Discord channels (Stable, PTB, Canary)
- 🧭 Discover Discord installs and suggested paths
- 🧩 Manage plugins and themes (list, install, update, remove)
- 🛒 Browse and search the BetterDiscord store
- 🖥️ Cross-platform support (Windows, macOS, Linux)
- 📦 Available via npm for easy distribution
- ⚡ Fast and lightweight Go binary

## Installation

### Via npm (Recommended)

```bash
npm install -g @betterdiscord/cli
```

### Via Go

```bash
go install github.com/betterdiscord/cli@latest
```

### Via winget (Windows)

```bash
winget install betterdiscord.cli
```

### Via Homebrew/Linuxbrew

```bash
brew install betterdiscord/tap/bdcli
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/BetterDiscord/cli/releases).

## Usage

### Global Options

```bash
bdcli --silent <command>   # Suppress non-error output
```

You can also set `BDCLI_SILENT=1` to silence output in automation.

### Install BetterDiscord

Install BetterDiscord to a specific Discord channel:

```bash
bdcli install --channel stable   # Install to Discord Stable
bdcli install --channel ptb      # Install to Discord PTB
bdcli install --channel canary   # Install to Discord Canary
```

Install BetterDiscord by providing a Discord install path:

```bash
bdcli install --path /path/to/Discord
```

### Uninstall BetterDiscord

Uninstall BetterDiscord from a specific Discord channel:

```bash
bdcli uninstall --channel stable   # Uninstall from Discord Stable
bdcli uninstall --channel ptb      # Uninstall from Discord PTB
bdcli uninstall --channel canary   # Uninstall from Discord Canary
```

Uninstall BetterDiscord by providing a Discord install path:

```bash
bdcli uninstall --path /path/to/Discord
```

Uninject BetterDiscord from all detected Discord installations (without deleting data):

```bash
bdcli uninstall --all
```

Fully uninstall BetterDiscord from all Discord installations and remove all BetterDiscord folders:

```bash
bdcli uninstall --full
```

### Check Version

```bash
bdcli version
```

### Update BetterDiscord

```bash
bdcli update
bdcli update --check
```

### Show BetterDiscord Info

```bash
bdcli info
```

### Discover Discord Installs

```bash
bdcli discover installs
bdcli discover paths
bdcli discover addons
```

### Manage Plugins

```bash
bdcli plugins list
bdcli plugins info <name>
bdcli plugins install <name|id|url>
bdcli plugins update <name|id|url>
bdcli plugins update <name|id> --check    # Check for updates without installing
bdcli plugins remove <name|id>
```

### Manage Themes

```bash
bdcli themes list
bdcli themes info <name>
bdcli themes install <name|id|url>
bdcli themes update <name|id|url>
bdcli themes update <name|id> --check     # Check for updates without installing
bdcli themes remove <name|id>
```

### Browse the Store

```bash
bdcli store search <query>
bdcli store show <id|name>

bdcli store plugins search <query>
bdcli store plugins show <id|name>

bdcli store themes search <query>
bdcli store themes show <id|name>
```

### Shell Completions

```bash
bdcli completion bash
bdcli completion zsh
bdcli completion fish
```

### Help

```bash
bdcli --help
bdcli [command] --help
```

### Automation

For scripts and CI jobs, you can suppress non-error output:

```bash
# One-off command
bdcli --silent install --channel stable

# Environment variable (applies to all commands)
BDCLI_SILENT=1 bdcli update
```

### CLI Help Output

```
A cross-platform CLI for installing, updating, and managing BetterDiscord.

Usage:
   bdcli [flags]
   bdcli [command]

Available Commands:
   completion  Generate shell completions
   discover    Discover Discord installations and related data
   help        Help about any command
   info        Displays information about BetterDiscord installation
   install     Installs BetterDiscord to your Discord
   plugins     Manage BetterDiscord plugins
   store       Browse and search the BetterDiscord store
   themes      Manage BetterDiscord themes
   uninstall   Uninstalls BetterDiscord from your Discord
   update      Update BetterDiscord to the latest version
   version     Print the version number

Flags:
       --silent   Suppress non-error output
   -h, --help     help for bdcli

Use "bdcli [command] --help" for more information about a command.
```

## Supported Platforms

- **Windows** (x64, ARM64, x86)
- **macOS** (x64, ARM64/M1/M2)
- **Linux** (x64, ARM64, ARM)

## Development

### Prerequisites

- [Go](https://go.dev/) 1.26 or higher
- [Task](https://taskfile.dev/) (optional, for task automation)
- [GoReleaser](https://goreleaser.com/) (for releases)

### Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/BetterDiscord/cli.git
cd cli
task setup  # Or: go mod download
```

### Available Tasks

Run `task --list-all` to see all available tasks:

```bash
# Development
task run             # Run the CLI (pass args with: task run -- install stable)

# Building
task build           # Build for current platform
task build:all       # Build for all platforms (GoReleaser)

# Testing
task test            # Run tests
task test:verbose    # Run tests with verbose output
task coverage        # Run tests with coverage summary
task coverage:html   # Generate HTML coverage report

# Code Quality
task fmt             # Format Go files
task vet             # Run go vet
task lint            # Run golangci-lint
task check           # Run fix, fmt, vet, lint, test

# Release
task release:snapshot # Test release build
task release          # Create release (requires tag)

# Cleaning
task clean           # Remove build and debug artifacts
```

### Running Locally

```bash
# Run directly
go run main.go install stable

# Or use Task
task run -- install stable
```

### Building

```bash
# Build for current platform
task build

# Build for all platforms
task build:all

# Output will be in ./dist/
```

### Testing

```bash
# Run all tests
task test

# Run with coverage
task coverage
```

### Releasing

1. Create and push a new tag:

   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```

2. GitHub Actions will automatically build and create a draft release

3. Edit the release notes and publish

4. Publish to npm:

   ```bash
   npm publish
   ```

## Project Structure

```py
.
├── cmd/                  # Cobra commands
│   ├── install.go       # Install command
│   ├── update.go        # Update command
│   ├── info.go          # Info command
│   ├── discover.go      # Discover command
│   ├── plugins.go       # Plugins commands
│   ├── themes.go        # Themes commands
│   ├── store.go         # Store commands
│   ├── uninstall.go     # Uninstall command
│   ├── version.go       # Version command
│   └── root.go          # Root command
├── internal/            # Internal packages
│   ├── betterdiscord/  # BetterDiscord installation logic
│   ├── discord/        # Discord path resolution and injection
│   ├── models/         # Data models
│   └── utils/          # Utility functions
├── main.go             # Entry point
├── Taskfile.yml        # Task automation
└── .goreleaser.yaml    # Release configuration
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Links

- [BetterDiscord Website](https://betterdiscord.app/)
- [BetterDiscord Documentation](https://docs.betterdiscord.app/)
- [Issue Tracker](https://github.com/BetterDiscord/cli/issues)
- [npm Package](https://www.npmjs.com/package/@betterdiscord/cli)

## Acknowledgments

Built with:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GoReleaser](https://goreleaser.com/) - Release automation
- [Task](https://taskfile.dev/) - Task runner

---

Made with ❤️ by the BetterDiscord Team
