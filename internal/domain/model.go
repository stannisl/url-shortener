package domain

import "time"

type RedirectModel struct {
	ShortUrl   string
	UserAgent  string
	RedirectAt time.Time
}

type ShortenUrlModel struct {
	OriginUrl string
	ShortUrl  string
}
