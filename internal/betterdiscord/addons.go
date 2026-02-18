package betterdiscord

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

type AddonKind string

const (
	AddonPlugin AddonKind = "plugin"
	AddonTheme  AddonKind = "theme"
)

var pluginExtensions = []string{".plugin.js"}
var themeExtensions = []string{".theme.css"}

type AddonEntry struct {
	BaseName     string    `json:"name"`
	FullFilename string    `json:"filename"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	Modified     time.Time `json:"modified"`
	Meta         Meta      `json:"meta"`
}

// ResolvedAddon holds both local and store metadata for an addon.
type ResolvedAddon struct {
	Store *models.StoreAddon // Metadata from store (nil if not found)
	URL   string             // Download URL
}

// ListAddons returns the locally installed addons for the given kind.
func ListAddons(kind AddonKind) ([]AddonEntry, error) {
	dir, err := addonDir(kind)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []AddonEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !isAddonFile(kind, e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		outerTrim := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		out = append(out, AddonEntry{
			BaseName:     strings.TrimSuffix(outerTrim, filepath.Ext(outerTrim)),
			FullFilename: e.Name(),
			Path:         filepath.Join(dir, e.Name()),
			Size:         info.Size(),
			Modified:     info.ModTime(),
			Meta:         parseJSDoc(string(contents)),
		})
	}
	return out, nil
}

// FindAddon searches for an installed addon by identifier (name, filename, or meta name).
// Returns the addon entry if found, or nil if not found.
func FindAddon(kind AddonKind, identifier string) *AddonEntry {
	items, err := ListAddons(kind)
	if err != nil {
		return nil
	}

	lower := strings.ToLower(identifier)
	for i := range items {
		// Match by filename (case-insensitive)
		if strings.ToLower(items[i].FullFilename) == lower || strings.ToLower(items[i].FullFilename) == lower+".plugin.js" || strings.ToLower(items[i].FullFilename) == lower+".theme.css" {
			return &items[i]
		}
		// Match by base name (case-insensitive)
		if strings.ToLower(items[i].BaseName) == lower {
			return &items[i]
		}
		// Match by meta name (case-insensitive)
		if strings.ToLower(items[i].Meta.Name) == lower {
			return &items[i]
		}
	}
	return nil
}

// InstallAddon installs an addon. Identifier can be:
// - A direct URL (https://...)
// - An addon ID (numeric)
// - An addon name (string)
// Returns the destination path and resolved addon metadata.
func InstallAddon(kind AddonKind, identifier string) (*ResolvedAddon, error) {
	dir, err := addonDir(kind)
	if err != nil {
		return nil, err
	}

	resolved := &ResolvedAddon{}

	// Case 1: Direct URL
	if utils.IsURL(identifier) {
		dest, err := downloadAddon(kind, dir, identifier)
		if err != nil {
			return nil, err
		}
		resolved.URL = dest
		return resolved, nil
	}

	// Case 2: ID or Name - query the store API
	addon, err := FetchAddonFromStore(identifier)
	if err != nil {
		return nil, fmt.Errorf("addon not found: %w", err)
	}

	resolved.Store = addon

	// Use the latest_source_url from the store if available
	downloadURL := addon.LatestSourceURL
	if downloadURL == "" {
		// Fallback: use the download redirect URL
		url, err := GetAddonDownloadURL(addon.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve download URL: %w", err)
		}
		downloadURL = url
	}

	dest, err := downloadAddon(kind, dir, downloadURL)
	if err != nil {
		return nil, err
	}

	resolved.URL = dest
	LogAddonInfo(addon)
	return resolved, nil
}

// RemoveAddon deletes a local addon by name or filename.
func RemoveAddon(kind AddonKind, identifier string) error {
	dir, err := addonDir(kind)
	if err != nil {
		return err
	}

	candidates := candidateFilenames(kind, identifier)
	for _, name := range candidates {
		full := filepath.Join(dir, name)
		if _, statErr := os.Stat(full); statErr == nil {
			return os.Remove(full)
		}
	}

	return fmt.Errorf("addon %s not found", identifier)
}

// UpdateAddon removes then installs again, returning resolved metadata.
func UpdateAddon(kind AddonKind, identifier string) (*ResolvedAddon, error) {
	_ = RemoveAddon(kind, identifier)
	return InstallAddon(kind, identifier)
}

func addonDir(kind AddonKind) (string, error) {
	inst := GetInstallation()
	// if err := inst.Prepare(); err != nil {
	// 	return "", err
	// }
	switch kind {
	case AddonPlugin:
		return inst.Plugins(), nil
	case AddonTheme:
		return inst.Themes(), nil
	default:
		return "", fmt.Errorf("unknown addon kind: %s", kind)
	}
}

func downloadAddon(kind AddonKind, dir, rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	base := filepath.Base(parsed.Path)
	if base == "." || base == "/" || base == "" {
		base = fmt.Sprintf("addon_%d", time.Now().Unix())
	}

	// Ensure correct extension for type when missing
	if !isAddonFile(kind, base) {
		switch kind {
		case AddonPlugin:
			base = base + ".plugin.js"
		case AddonTheme:
			base = base + ".theme.css"
		}
	}

	dest := filepath.Join(dir, base)
	if _, err := utils.DownloadFile(rawURL, dest); err != nil {
		return "", err
	}
	return dest, nil
}

func isAddonFile(kind AddonKind, name string) bool {
	lower := strings.ToLower(name)
	switch kind {
	case AddonPlugin:
		for _, ext := range pluginExtensions {
			if strings.HasSuffix(lower, ext) {
				return true
			}
		}
	case AddonTheme:
		for _, ext := range themeExtensions {
			if strings.HasSuffix(lower, ext) {
				return true
			}
		}
	}
	return false
}

func candidateFilenames(kind AddonKind, identifier string) []string {
	var out []string
	lower := strings.ToLower(identifier)
	if isAddonFile(kind, lower) {
		out = append(out, identifier)
	}
	switch kind {
	case AddonPlugin:
		for _, ext := range pluginExtensions {
			out = append(out, lower+ext)
		}
	case AddonTheme:
		for _, ext := range themeExtensions {
			out = append(out, lower+ext)
		}
	}
	return out
}

// LogLocalAddonInfo prints detailed information about a locally installed addon.
func LogLocalAddonInfo(entry *AddonEntry) {
	name := entry.Meta.Name
	if name == "" {
		name = entry.BaseName
	}

	// Header with name and version
	versionStr := ""
	if entry.Meta.Version != "" {
		versionStr = fmt.Sprintf(" v%s", entry.Meta.Version)
	}
	output.Printf("📦 %s%s\n", name, versionStr)

	// Author with link if available
	if entry.Meta.Author != "" {
		authorStr := entry.Meta.Author
		if entry.Meta.AuthorLink != "" {
			authorStr = fmt.Sprintf("%s (%s)", authorStr, entry.Meta.AuthorLink)
		}
		output.Printf("   By: %s\n", authorStr)
	}

	// Description
	if entry.Meta.Description != "" {
		output.Printf("   %s\n", entry.Meta.Description)
	}

	// File information section
	output.Blank()
	fileName := entry.FullFilename
	if fileName == "" {
		fileName = entry.BaseName
	}
	output.Printf("   📁 File: %s\n", fileName)
	output.Printf("   💾 Size: %.1f KB\n", float64(entry.Size)/1024.0)
	output.Printf("   🕐 Modified: %s\n", entry.Modified.Format(output.DateTimeFormat))

	// Links section
	hasLinks := entry.Meta.Website != "" || entry.Meta.Source != "" || entry.Meta.AuthorLink != "" || entry.Meta.Donate != "" || entry.Meta.Patreon != "" || entry.Meta.Invite != ""

	if hasLinks {
		output.Blank()
	}

	if entry.Meta.Website != "" {
		output.Printf("   🌐 Website: %s\n", entry.Meta.Website)
	}
	if entry.Meta.Source != "" {
		output.Printf("   🔗 Source: %s\n", entry.Meta.Source)
	}
	if entry.Meta.Donate != "" {
		output.Printf("   💜 Donate: %s\n", entry.Meta.Donate)
	}
	if entry.Meta.Patreon != "" {
		output.Printf("   💰 Patreon: %s\n", entry.Meta.Patreon)
	}
	if entry.Meta.Invite != "" {
		output.Printf("   💬 Discord: https://discord.gg/%s\n", entry.Meta.Invite)
	}
	output.Blank()
}
