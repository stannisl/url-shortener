package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/lib/pq"
	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/dbpg"
)

type UrlRepository struct {
	db *dbpg.DB
}

func NewUrlRepository(db *dbpg.DB) *UrlRepository {
	return &UrlRepository{db: db}
}

func (r *UrlRepository) GetOriginUrl(ctx context.Context, shortUrl string) (*domain.ShortenUrlModel, error) {
	query := `SELECT id, original_url, short_code 
				FROM urls 
			    WHERE short_code = $1`

	var shortenUrlModel domain.ShortenUrlModel

	err := r.db.QueryRowContext(ctx, query, shortUrl).Scan(
		&shortenUrlModel.ID,
		&shortenUrlModel.OriginalUrl,
		&shortenUrlModel.ShortCode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrModelNotFound
		}
		return nil, err
	}

	return &shortenUrlModel, nil
}

func (r *UrlRepository) CreateShortUrl(ctx context.Context, originUrl string) (*domain.ShortenUrlModel, error) {
	query := `
			INSERT INTO urls (original_url, short_code) 
			VALUES ($1, $2) 
			RETURNING (original_url, short_code)
			`

	var result domain.ShortenUrlModel

	var shortUrl string
	retries := 5
	var pgErr *pq.Error
	for i := 0; i < retries; i++ {
		shortUrl = generateRandomUrl(min(len(originUrl), 10))
		err := r.db.QueryRowContext(ctx, query, originUrl, shortUrl).Scan(&result.OriginalUrl, &result.ShortCode)

		if err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					continue
				}
				return nil, err
			}
		}
		return &domain.ShortenUrlModel{
			OriginalUrl: originUrl,
			ShortCode:   shortUrl,
		}, nil
	}

	if pgErr.Code == "23505" {
		return nil, domain.ErrNotUnique
	}

	return &domain.ShortenUrlModel{
		OriginalUrl: originUrl,
		ShortCode:   shortUrl,
	}, nil
}

func generateRandomUrl(length int) string {
	randomBytes := make([]byte, length)
	_, _ = rand.Read(randomBytes)

	url := base64.URLEncoding.EncodeToString(randomBytes)

	return url[:length]
}
