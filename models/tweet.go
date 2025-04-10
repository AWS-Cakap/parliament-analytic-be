package models

import "time"

type Tweet struct {
	ID        uint
	Username  string
	Text      string
	CreatedAt time.Time
	Sentiment string
}
