package postgres

import (
	"github.com/stannisl/url-shortener/internal/config"
	"github.com/wb-go/wbf/dbpg"
)

func NewDB(config config.PostgresConfig) (*dbpg.DB, error) {
	db, err := dbpg.New(config.DatabaseUrl, nil, &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
