package service

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stannisl/url-shortener/internal/repository"
)

type UrlService struct {
	urlRepository       UrlRepository
	analyticsRepository repository.AnalyticsRepository
}

func NewUrlService(urlRepository UrlRepository) *UrlService {
	return &UrlService{
		urlRepository: urlRepository,
	}
}

func (us *UrlService) CreateShortenUrl(ctx context.Context, originUrl string) (*domain.ShortenUrlModel, error) {
	shortenURLS, err := us.urlRepository.CreateShortUrl(ctx, originUrl)
	if err != nil {
		return nil, err
	}

	return shortenURLS, nil
}

func (us *UrlService) GetOriginUrl(ctx context.Context, shortUrl string) (*domain.ShortenUrlModel, error) {
	shortenUrls, err := us.urlRepository.GetOriginUrl(ctx, shortUrl)
	if err != nil {
		return nil, err
	}

	return shortenUrls, nil
}
