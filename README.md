# BetterDiscord CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/BetterDiscord/cli)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/BetterDiscord/cli)](https://github.com/BetterDiscord/cli/releases)
[![License](https://img.shields.io/github/license/BetterDiscord/cli)](LICENSE)
[![npm](https://img.shields.io/npm/v/@betterdiscord/cli)](https://www.npmjs.com/package/@betterdiscord/cli)

A cross-platform command-line interface for installing, updating, and managing [BetterDiscord](https://betterdiscord.app/).

## Features

- 🚀 Easy installation and uninstallation of BetterDiscord
- 🔄 Support for multiple Discord channels (Stable, PTB, Canary)
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

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/BetterDiscord/cli/releases).

## Usage

### Install BetterDiscord

Install BetterDiscord to a specific Discord channel:

```bash
bdcli install stable   # Install to Discord Stable
bdcli install ptb      # Install to Discord PTB
bdcli install canary   # Install to Discord Canary
```

### Uninstall BetterDiscord

Uninstall BetterDiscord from a specific Discord channel:

```bash
bdcli uninstall stable   # Uninstall from Discord Stable
bdcli uninstall ptb      # Uninstall from Discord PTB
bdcli uninstall canary   # Uninstall from Discord Canary
```

### Check Version

```bash
bdcli version
```

### Help

```bash
bdcli --help
bdcli <command> --help
```

## Supported Platforms

- **Windows** (x64, ARM64, x86)
- **macOS** (x64, ARM64/M1/M2)
- **Linux** (x64, ARM64, ARM)

## Development

### Prerequisites

- [Go](https://go.dev/) 1.19 or higher
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

Run `task --list` to see all available tasks:

```bash
# Development
task run              # Run the CLI
task run:install      # Test install command
task run:uninstall    # Test uninstall command

# Building
task build            # Build for current platform
task build:all        # Build for all platforms
task install          # Install to $GOPATH/bin

# Testing
task test             # Run tests
task test:coverage    # Run tests with coverage

# Code Quality
task lint             # Run linter
task fmt              # Format code
task vet              # Run go vet
task check            # Run all checks

# Release
task release:snapshot # Test release build
task release          # Create release (requires tag)
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
task test:coverage
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
