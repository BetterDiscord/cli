package utils

import "net/url"

// IsURL checks if a string is a valid URL
func IsURL(input string) bool {
	parsed, err := url.Parse(input)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}
