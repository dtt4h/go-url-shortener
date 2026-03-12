package model

import "time"

type URLEvent struct {
	EventType   string
	ShortCode   string
	OriginalURL string
	Timestamp   time.Time
}
