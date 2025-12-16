package service

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
)

type AnalyticsRepository interface {
	GetSummary(ctx context.Context, urlID int) (*domain.AnalyticsSummary, error)
	AddAnalytics(ctx context.Context, redirectModel domain.UrlAnalyticsModel) error
}

type UrlRepository interface {
	GetOriginUrl(ctx context.Context, shortUrl string) (*domain.ShortenUrlModel, error)
	CreateShortUrl(ctx context.Context, originUrl string) (*domain.ShortenUrlModel, error)
}
