// BetterDiscord's Injection Script
const path = require("path");
const electron = require("electron");

// Windows and macOS both use the fixed global BetterDiscord folder but
// Electron gives the postfixed version of userData, so go up a directory
let userConfig = path.join(electron.app.getPath("userData"), "..");

// If we're on Linux there are a couple cases to deal with
if (process.platform !== "win32" && process.platform !== "darwin") {
    // Use || instead of ?? because a falsey value of "" is invalid per XDG spec
    userConfig = process.env.XDG_CONFIG_HOME || path.join(process.env.HOME, ".config");

    // HOST_XDG_CONFIG_HOME is set by flatpak, so use without validation if set
    if (process.env.HOST_XDG_CONFIG_HOME) userConfig = process.env.HOST_XDG_CONFIG_HOME;
}

require(path.join(userConfig, "BetterDiscord", "data", "betterdiscord.asar"));

// Discord's Default Export
module.exports = require("./core.asar");