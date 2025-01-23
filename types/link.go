package types

import "time"

type Link struct {
	Id       string    `json:"id"`
	Short    string    `json:"short"`
	Long     float64   `json:"long"`
	CreateAt time.Time `json:"create_at"`
}
