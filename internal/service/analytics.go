package service

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
)

type AnalyticsService struct {
	analyticsRepository AnalyticsRepository
	urlRepository       UrlRepository
}

func NewAnalyticsService(
	urlRepository UrlRepository,
	analyticsRepository AnalyticsRepository,
) *AnalyticsService {
	return &AnalyticsService{
		urlRepository:       urlRepository,
		analyticsRepository: analyticsRepository,
	}
}

func (as *AnalyticsService) HandleRedirect(ctx context.Context, shortUrl string, redirectInfo domain.UrlAnalyticsModel) error {
	shortenUrl, err := as.urlRepository.GetOriginUrl(ctx, shortUrl)
	if err != nil {
		return err
	}

	redirectInfo.UrlId = shortenUrl.ID

	err = as.analyticsRepository.AddAnalytics(ctx, redirectInfo)
	if err != nil {
		return err
	}

	return nil
}

func (as *AnalyticsService) GetStats(ctx context.Context, shortUrl string) (*domain.AnalyticsSummary, error) {
	shortenUrl, err := as.urlRepository.GetOriginUrl(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	analytics, err := as.analyticsRepository.GetSummary(ctx, shortenUrl.ID)
	if err != nil {
		return nil, err
	}
	return analytics, nil
}
