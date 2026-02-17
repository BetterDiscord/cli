package models

import "time"

type StoreAddon struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	FileName           string    `json:"file_name"`
	Type               string    `json:"type"` // "plugin" or "theme"
	Description        string    `json:"description"`
	Version            string    `json:"version"`
	Author             Author    `json:"author"`
	Likes              int       `json:"likes"`
	Downloads          int       `json:"downloads"`
	Tags               []string  `json:"tags"`
	ThumbnailURL       string    `json:"thumbnail_url"`
	LatestSourceURL    string    `json:"latest_source_url"`
	InitialReleaseDate time.Time `json:"initial_release_date"`
	LatestReleaseDate  time.Time `json:"latest_release_date"`
	Guild              *Guild    `json:"guild"`
}

type Author struct {
	GitHubID          string `json:"github_id"`
	GitHubName        string `json:"github_name"`
	DisplayName       string `json:"display_name"`
	DiscordName       string `json:"discord_name"`
	DiscordAvatarHash string `json:"discord_avatar_hash"`
	DiscordSnowflake  string `json:"discord_snowflake"`
	Guild             *Guild `json:"guild"`
}

type Guild struct {
	Name       string `json:"name"`
	Snowflake  string `json:"snowflake"`
	InviteLink string `json:"invite_link"`
	AvatarHash string `json:"avatar_hash"`
}
