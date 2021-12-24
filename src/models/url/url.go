package model

import "time"

type Url struct {
	CreatedAt   time.Time
	RedirectUrl string
	ShortenedId string
}
