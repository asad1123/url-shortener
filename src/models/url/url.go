package model

import "time"

type Url struct {
	CreatedAt   time.Time
	RedirectUrl string
	ShortenedId string
}

type UrlUsage struct {
	ShortenedId string
	AccessedAt  time.Time
}

type UrlUsageRequestSchema struct {
	// this will be of the form <number><time period>
	// e.g. 24h -> 24 hours
	Since string `form:"since"`
}

type UrlUsageResponseSchema struct {
	ShortenedId string
	Count       int
}
