package betterdiscord

import (
	"log"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

func (i *BDInstall) download() error {
	if i.hasDownloaded {
		log.Printf("✅ Already downloaded to %s", i.asar)
		return nil
	}

	resp, err := utils.DownloadFile("https://betterdiscord.app/Download/betterdiscord.asar", i.asar)
	if err == nil {
		version := resp.Header.Get("x-bd-version")
		log.Printf("✅ Downloaded BetterDiscord version %s from the official website", version)
		i.hasDownloaded = true
		return nil
	} else {
		log.Printf("❌ Failed to download BetterDiscord from official website")
		log.Printf("❌ %s", err.Error())
		log.Printf("")
		log.Printf("#### Falling back to GitHub...")
	}

	// Get download URL from GitHub API
	apiData, err := utils.DownloadJSON[models.GitHubRelease]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
	if err != nil {
		log.Printf("❌ Failed to get asset url from GitHub")
		log.Printf("❌ %s", err.Error())
		return err
	}

	var index = 0
	for i, asset := range apiData.Assets {
		if asset.Name == "betterdiscord.asar" {
			index = i
			break
		}
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		log.Printf("✅ Found BetterDiscord: %s", downloadUrl)
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, i.asar)
	if err != nil {
		log.Printf("❌ Failed to download BetterDiscord from GitHub")
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Downloaded BetterDiscord version %s from GitHub", version)
	i.hasDownloaded = true

	return nil
}
