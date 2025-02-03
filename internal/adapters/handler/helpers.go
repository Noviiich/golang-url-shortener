package handler

import (
	"net/url"
	"regexp"
)

func IsValidLink(u string) bool {
	re := regexp.MustCompile(`^(http|https)://`)
	if !re.MatchString(u) {
		return false
	}

	parsedURL, err := url.ParseRequestURI(u)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}
