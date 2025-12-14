package service

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stannisl/url-shortener/internal/repository"
)

type AnalyticsService struct {
	analyticsRepository repository.AnalyticsRepository
	urlRepository       repository.UrlRepository
}

func NewAnalyticsService(
	urlRepository repository.UrlRepository,
	analyticsRepository repository.AnalyticsRepository,
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
