<h1 align="center">BetterDiscord Command-Line-Interface Installer</h1>

<p align="center">
  <a href="#overview">Overview</a> |
  <a href="#development">Development</a> |
  <a href="#contributors">Contributors</a>
</p>

<p align="center">
  <br/>
  A simple standalone command-line-interface which automates the installation, removal and maintenance of <a href="https://github.com/BetterDiscord/BetterDiscord">BetterDiscord</a>.
  <br/>
  <br/>
  <a href="https://betterdiscord.app/invite" target="_blank">
    <img src="https://img.shields.io/badge/discord-join-green?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0OCIgaGVpZ2h0PSI0OCIgdmlld0JveD0iMCAwIDQ4IDQ4IiBmaWxsPSJub25lIj4NCjxwYXRoIGQ9Ik0xNi41MzUzIDUuNDUwNTNDMzMuODIzMyAtMS40NjgwOSA1MC44ODE1IDE2Ljg4NjMgNDEuNTkyNSAzMy42MDY2QzM3LjM3MzIgNDEuMjAxMiAyNi44OTA0IDQ3LjMxNyAxNC42ODQyIDQxLjUzMjZMNi4xOTk5IDQzLjk1NjdDNC44ODM0NiA0NC4zMzI4IDMuNjY2OTMgNDMuMTIxMiA0LjAzNjM0IDQxLjgwMzdDNC41NDI3MiAzOS45OTc2IDUuNzQyNTcgMzUuNzM5OCA2LjQ0NDEgMzMuNDMyM0MxLjE4Mjc5IDI0LjA0NCA0LjczMDUgMTAuMTc0OCAxNi41MzUzIDUuNDUwNTNaTTE1Ljk5NTQgMjAuMjQ5NkMxNS45OTU0IDIwLjkzOTkgMTYuNTU1IDIxLjQ5OTYgMTcuMjQ1NCAyMS40OTk2SDMwLjc0OThDMzEuNDQwMSAyMS40OTk2IDMxLjk5OTggMjAuOTM5OSAzMS45OTk4IDIwLjI0OTZDMzEuOTk5OCAxOS41NTkyIDMxLjQ0MDEgMTguOTk5NiAzMC43NDk4IDE4Ljk5OTZIMTcuMjQ1NEMxNi41NTUgMTguOTk5NiAxNS45OTU0IDE5LjU1OTIgMTUuOTk1NCAyMC4yNDk2Wk0xNy4yNDk4IDI2LjQ3NDZDMTYuNTU5NCAyNi40NzQ2IDE1Ljk5OTggMjcuMDM0MiAxNS45OTk4IDI3LjcyNDZDMTUuOTk5OCAyOC40MTQ5IDE2LjU1OTQgMjguOTc0NiAxNy4yNDk4IDI4Ljk3NDZIMjYuNzQ5OEMyNy40NDAxIDI4Ljk3NDYgMjcuOTk5OCAyOC40MTQ5IDI3Ljk5OTggMjcuNzI0NkMyNy45OTk4IDI3LjAzNDIgMjcuNDQwMSAyNi40NzQ2IDI2Ljc0OTggMjYuNDc0NkgxNy4yNDk4WiIgZmlsbD0iIzNhNzFjMSIvPg0KPC9zdmc+" alt="Chat"/>
  </a>
  <a href="https://github.com/BetterDiscord/Installer/blob/main/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/BetterDiscord/cli?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEwLjk2ODQgMi4zMjQ2NUMxMS41ODMgMS44NzYxNiAxMi40MTcgMS44NzYxNiAxMy4wMzE2IDIuMzI0NjVMMjAuNDUzNCA3Ljc0MDZDMjEuNDI5OSA4LjQ1MzE1IDIwLjkyNjggOS45OTgzNSAxOS43MTg5IDEwLjAwMDNINC4yODEwOEMzLjA3MzE4IDkuOTk4MzUgMi41NzAxMSA4LjQ1MzE1IDMuNTQ2NTcgNy43NDA2TDEwLjk2ODQgMi4zMjQ2NVpNMTMgNi4yNTAzNEMxMyA1LjY5ODA1IDEyLjU1MjMgNS4yNTAzNCAxMiA1LjI1MDM0QzExLjQ0NzcgNS4yNTAzNCAxMSA1LjY5ODA1IDExIDYuMjUwMzRDMTEgNi44MDI2MiAxMS40NDc3IDcuMjUwMzQgMTIgNy4yNTAzNEMxMi41NTIzIDcuMjUwMzQgMTMgNi44MDI2MiAxMyA2LjI1MDM0WiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTEuMjUgMTYuMDAwM0g5LjI1VjExLjAwMDNIMTEuMjVWMTYuMDAwM1oiIGZpbGw9IiMzYTcxYzEiLz4KPHBhdGggZD0iTTE0Ljc1IDE2LjAwMDNIMTIuNzVWMTEuMDAwM0gxNC43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNSAxNi4wMDAzSDE2LjI1VjExLjAwMDNIMTguNVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNzUgMTcuMDAwM0g1LjI1QzQuMDA3MzYgMTcuMDAwMyAzIDE4LjAwNzcgMyAxOS4yNTAzVjE5Ljc1MDNDMyAyMC4xNjQ1IDMuMzM1NzkgMjAuNTAwMyAzLjc1IDIwLjUwMDNIMjAuMjVDMjAuNjY0MiAyMC41MDAzIDIxIDIwLjE2NDUgMjEgMTkuNzUwM1YxOS4yNTAzQzIxIDE4LjAwNzcgMTkuOTkyNiAxNy4wMDAzIDE4Ljc1IDE3LjAwMDNaIiBmaWxsPSIjM2E3MWMxIi8+CjxwYXRoIGQ9Ik03Ljc1IDE2LjAwMDNINS41VjExLjAwMDNINy43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8L3N2Zz4K" alt="License"/>
  </a>
</p>

---

# Overview

This repository contains the source code for the BetterDiscord cli installer.

## Download

```ps
sudo npm @betterdiscord/cli
```

---

# Development

> This is a tutorial designed for people looking to contribute to, or work directly with the installer's source code. If you are just looking to download and install BetterDiscord, visit the [releases](https://github.com/BetterDiscord/installer/releases/) page of this repository.

## Prerequisites
- [npm](https://www.npmjs.com/)
- Command line of your choice.

## Building

### 1: Clone the repository.
```ps
git clone https://github.com/BetterDiscord/cli && cd cli
```
This will create a local copy of this repository and navigate you to the root folder of the repository.

### 2: Install Dependencies
Run this command at the root folder to install dependencies:
```ps
sudo npm install
```

---

# Contributors

For information on contributing to this project, please see [CONTRIBUTING.md](/CONTRIBUTING.md).

<a href="https://github.com/betterdiscord/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=betterdiscord/cli" />
</a>
