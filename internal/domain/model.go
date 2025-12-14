package domain

import "time"

type UrlAnalyticsModel struct {
	ID         int       `db:"id"`
	UrlId      int       `db:"url_id"`
	UserAgent  string    `db:"user_agent"`
	IPAddress  string    `db:"ip_address"`
	AccessedAt time.Time `db:"accessed_at"`
}

type ShortenUrlModel struct {
	ID          int       `db:"id"`
	OriginalUrl string    `db:"original_url"`
	ShortCode   string    `db:"short_code"`
	CreatedAt   time.Time `db:"created_at"`
}

type AnalyticsSummary struct {
	TotalClicks    int `db:"total_clicks"`
	UniqueVisitors int `db:"unique_visitors"`
}
