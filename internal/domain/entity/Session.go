package entity

import "time"

type Session struct {
	UUID   string    `json:"uuid"`
	Expire time.Time `json:"expire"`
	ID     int       `json:"id"`
}
