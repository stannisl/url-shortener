package repository

import (
	"context"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/dbpg"
)

type AnalyticsRepository struct {
	db *dbpg.DB
}

func NewAnalyticsRepository(db *dbpg.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context, urlID int) (*domain.AnalyticsSummary, error) {
	query := `SELECT 
    			count(*) as total_clicks,
				count(DISTINCT ip_address) as unique_clicks 
			  FROM url_analytics 
			  WHERE url_id = $1`

	var analytics domain.AnalyticsSummary
	err := r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.TotalClicks, &analytics.UniqueVisitors)
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

func (r *AnalyticsRepository) AddAnalytics(ctx context.Context, redirectModel domain.UrlAnalyticsModel) error {
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
