package betterdiscord

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

// FetchAddonFromStore queries the BetterDiscord Store API by name or ID.
// Returns addon metadata including download URL.
func FetchAddonFromStore(identifier string) (*models.StoreAddon, error) {
	apiURL := fmt.Sprintf("https://api.betterdiscord.app/v3/store/%s", identifier)

	addon, err := utils.DownloadJSON[models.StoreAddon](apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addon '%s' from store: %w", identifier, err)
	}

	return &addon, nil
}

// GetAddonDownloadURL resolves the final download URL for an addon by ID.
// It follows redirects from the BetterDiscord download page.
func GetAddonDownloadURL(id int) (string, error) {
	downloadURL := fmt.Sprintf("https://betterdiscord.app/gh-redirect?id=%d", id)

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "BetterDiscord/cli")

	// Create client that follows redirects and returns the final URL
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download page returned status %d", resp.StatusCode)
	}

	// The final URL after redirects is the download URL
	return resp.Request.URL.String(), nil
}

// LogAddonInfo prints addon information for the user.
func LogAddonInfo(addon *models.StoreAddon) {
	log.Printf("📦 Addon: %s", addon.Name)
	log.Printf("   Version: %s", addon.Version)
	log.Printf("   Author: %s", addon.Author.DisplayName)
	log.Printf("   Description: %s", addon.Description)
	log.Printf("   Downloads: %d | Likes: %d", addon.Downloads, addon.Likes)
	if len(addon.Tags) > 0 {
		log.Printf("   Tags: %v", addon.Tags)
	}
}

// ResolveAddonIdentifier attempts to parse identifier as int (ID) or string (name).
// Returns (id, name, isID) where isID indicates whether it was parsed as an ID.
func ResolveAddonIdentifier(identifier string) (int, string, bool) {
	if id, err := strconv.Atoi(identifier); err == nil {
		return id, "", true
	}
	return 0, identifier, false
}
