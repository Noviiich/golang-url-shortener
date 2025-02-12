package domain

import "time"

type Platform int

const (
	PlatformUnknown Platform = iota
	PlatformInstagram
	PlatformTwitter
	PlatformYouTube
)

func (p Platform) String() string {
	switch p {
	case PlatformInstagram:
		return "Instagram"
	case PlatformTwitter:
		return "Twitter"
	case PlatformYouTube:
		return "YouTube"
	default:
		return "Unknown"
	}
}

type Stats struct {
	Id        string    `json:"id" bson:"id"`
	Platform  Platform  `json:"platform" bson:"platform"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
