package repository

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/dbpg"
)

type AnalyticsRepository interface {
	GetAnalytics(ctx context.Context, originUrl string) (string, error)
	AddAnalytics(ctx context.Context, redirectModel domain.RedirectModel) error
}

type analyticsRepository struct {
	db *dbpg.DB
}

func NewAnalyticsRepository(db *dbpg.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetAnalytics(ctx context.Context, originUrl string) (string, error) {}

func (r *analyticsRepository) AddAnalytics(ctx context.Context, redirectModel domain.RedirectModel) error {
}
