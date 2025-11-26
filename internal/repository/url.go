package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/dbpg"
)

type UrlRepository interface {
	GetOriginUrl(ctx context.Context, shortUrl string) (*domain.ShortenUrlModel, error)
	CreateShortUrl(ctx context.Context, originUrl string) (*domain.ShortenUrlModel, error)
}

type urlRepository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) UrlRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) GetOriginUrl(ctx context.Context, shortUrl string) (*domain.ShortenUrlModel, error) {
	query := `SELECT origin_url FROM urls WHERE short_url = $1`

	var originUrl string

	err := r.db.QueryRowContext(ctx, query, shortUrl).Scan(&originUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ModelNotFoundErr
		}
		return nil, err
	}

	return &domain.ShortenUrlModel{
		OriginUrl: originUrl,
		ShortUrl:  shortUrl,
	}, nil
}

func (r *urlRepository) CreateShortUrl(ctx context.Context, originUrl string) (*domain.ShortenUrlModel, error) {

}
