package domain

import "time"

type Link struct {
	Id          string    `bson:"id" json:"id"`
	OriginalURL string    `bson:"original_url" json:"original_url"`
	CreateAt    time.Time `bson:"create_at" json:"create_at"`
	Stats       []Stats   `bson:"-" json:"stats"`
}
