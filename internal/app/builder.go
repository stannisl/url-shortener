package app

import (
	"errors"
	"fmt"

	"github.com/stannisl/url-shortener/internal/config"
	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/repository"
	"github.com/stannisl/url-shortener/internal/service"
	"github.com/stannisl/url-shortener/internal/transport/http/router"
	"github.com/stannisl/url-shortener/internal/transport/http/server"
	"github.com/wb-go/wbf/dbpg"
)

type Builder struct {
	config *config.Config
	db     *dbpg.DB
	logger logger.Logger
	router router.Router
	err    error
}

func NewBuilder() *Builder {
	return &Builder{err: nil}
}

func (b *Builder) WithConfig(c *config.Config) *Builder {
	if c == nil {
		b.err = errors.Join(b.err, errors.New("config is nil"))
	}

	b.config = c
	return b
}

func (b *Builder) WithLogger(l logger.Logger) *Builder {
	if l == nil {
		b.err = errors.Join(b.err, errors.New("logger is nil"))
	}

	b.logger = l
	return b
}

func (b *Builder) WithPostgresDB(db *dbpg.DB) *Builder {
	if db == nil {
		b.err = errors.Join(b.err, errors.New("db is nil"))
	}

	b.db = db
	return b
}

func (b *Builder) Build() (*App, error) {
	if b.err != nil {
		return nil, fmt.Errorf("app builder errors: %w", b.err)
	}

	urlRepository := repository.NewUrlRepository(b.db)
	analyticsRepository := repository.NewAnalyticsRepository(b.db)

	urlService := service.NewUrlService(urlRepository)
	analyticsService := service.NewAnalyticsService(urlRepository, analyticsRepository)

	b.router = router.New(
		b.logger,
		urlService,
		analyticsService,
		b.config.ServerConfig.Mode,
	)

	httpServer, err := server.NewBuilder().
		WithPort(b.config.ServerConfig.Port).
		WithHost(b.config.ServerConfig.Host).
		WithHandler(b.router).
		Build()
	if err != nil {
		return nil, fmt.Errorf("server builder error: %w", err)
	}

	return &App{
		db:     b.db,
		config: b.config,
		logger: b.logger,
		server: httpServer,
		router: b.router,
	}, nil
}
