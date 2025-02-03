package model

import "time"

type Link struct {
	ShortId     string    `bson:"short_id" json:"short_id"`
	OriginalURL string    `bson:"original_url" json:"original_url"`
	CreateAt    time.Time `bson:"create_at" json:"create_at"`
	Clicks      int       `bson:"clicks" json:"clicks"`
}
