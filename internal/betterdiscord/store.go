package betterdiscord

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
func GetAddonDownloadURL(id int) (s string, err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download page returned status %d", resp.StatusCode)
	}

	// The final URL after redirects is the download URL
	return resp.Request.URL.String(), nil
}

// LogAddonInfo prints detailed addon information for the user.
func LogAddonInfo(addon *models.StoreAddon) {
	// Header with name, version, and type
	var typeStr string
	if addon.Type != "" {
		typeStr = fmt.Sprintf(" [%s]", strings.ToUpper(addon.Type))
	}
	log.Printf("📦 %s v%s%s", addon.Name, addon.Version, typeStr)

	// Author with GitHub link
	authorStr := addon.Author.DisplayName
	if addon.Author.GitHubName != "" {
		authorStr = fmt.Sprintf("%s (github.com/%s)", authorStr, addon.Author.GitHubName)
	}
	log.Printf("   By: %s", authorStr)

	// Description
	if addon.Description != "" {
		log.Printf("   %s", addon.Description)
	}

	// Stats line
	log.Printf("")
	log.Printf("   📊 Downloads: %d  |  👍 Likes: %d", addon.Downloads, addon.Likes)

	// Tags
	if len(addon.Tags) > 0 {
		tagsStr := strings.Join(addon.Tags, ", ")
		log.Printf("   🏷️  Tags: %s", tagsStr)
	}

	// Release dates
	log.Printf("")
	if !addon.InitialReleaseDate.IsZero() {
		log.Printf("   📅 Released: %s", addon.InitialReleaseDate.Format("Jan 2, 2006"))
	}
	if !addon.LatestReleaseDate.IsZero() {
		log.Printf("   🔄 Updated: %s", addon.LatestReleaseDate.Format("Jan 2, 2006"))
	}

	// Links
	if addon.LatestSourceURL != "" {
		log.Printf("   🔗 Source: %s", addon.LatestSourceURL)
	}
	if addon.Guild != nil && addon.Guild.InviteLink != "" {
		log.Printf("   💬 Server: %s", addon.Guild.InviteLink)
	}
	log.Printf("")
}

// ResolveAddonIdentifier attempts to parse identifier as int (ID) or string (name).
// Returns (id, name, isID) where isID indicates whether it was parsed as an ID.
func ResolveAddonIdentifier(identifier string) (int, string, bool) {
	if id, err := strconv.Atoi(identifier); err == nil {
		return id, "", true
	}
	return 0, identifier, false
}

// FetchAddonsOfType fetches all addons of a specific type from the store.
// Kind can be "plugin", "theme", or "addon" for all types.
func FetchAddonsOfType(kind string) ([]models.StoreAddon, error) {
	endpoint := kind
	if kind == "" {
		endpoint = "addons"
	}
	apiURL := fmt.Sprintf("https://api.betterdiscord.app/v3/store/%s", endpoint)

	addons, err := utils.DownloadJSON[[]models.StoreAddon](apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s from store: %w", endpoint, err)
	}

	return addons, nil
}

// SearchAddons performs a client-side search on addon slice.
// Searches addon Name, Description, Author DisplayName, and FileName.
func SearchAddons(addons []models.StoreAddon, query string) []models.StoreAddon {
	if query == "" {
		return addons
	}

	query = strings.ToLower(query)
	var results []models.StoreAddon

	for _, addon := range addons {
		if strings.Contains(strings.ToLower(addon.Name), query) ||
			strings.Contains(strings.ToLower(addon.Description), query) ||
			strings.Contains(strings.ToLower(addon.Author.DisplayName), query) ||
			strings.Contains(strings.ToLower(addon.FileName), query) {
			results = append(results, addon)
		}
	}

	return results
}
