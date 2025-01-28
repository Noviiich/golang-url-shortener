package types

import "time"

type Link struct {
	// Id       string    `bson:"_id,omitempty" json:"id"`
	Short    string    `bson:"short" json:"short"`
	Long     string    `bson:"long" json:"long"`
	CreateAt time.Time `bson:"create_at" json:"create_at"`
}
