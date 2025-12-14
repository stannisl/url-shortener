package repository

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/dbpg"
)

type AnalyticsRepository interface {
	GetSummary(ctx context.Context, urlID int) (*domain.AnalyticsSummary, error)
	AddAnalytics(ctx context.Context, redirectModel domain.UrlAnalyticsModel) error
}

type analyticsRepository struct {
	db *dbpg.DB
}

func NewAnalyticsRepository(db *dbpg.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetSummary(ctx context.Context, urlID int) (*domain.AnalyticsSummary, error) {
	query := `SELECT 
    			count(*) as total_clicks,
				count(DISTINCT ip_address) as unique_clicks 
			  FROM url_analytics 
			  WHERE url_id = $1
	`

	var analytics domain.AnalyticsSummary
	err := r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.TotalClicks, &analytics.UniqueVisitors)
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

func (r *analyticsRepository) AddAnalytics(ctx context.Context, redirectModel domain.UrlAnalyticsModel) error {
	query := `INSERT INTO url_analytics (url_id, user_agent, ip_address)
			  VALUES ($1, $2, $3)
			  RETURNING (url_id, user_agent, ip_address)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		redirectModel.UrlId,
		redirectModel.UserAgent,
		redirectModel.IPAddress,
	)

	return err
}
