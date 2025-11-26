package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/stannisl/url-shortener/internal/config"
	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/transport/http/router"
	"github.com/stannisl/url-shortener/internal/transport/http/server"
	"github.com/wb-go/wbf/dbpg"
	"golang.org/x/sync/errgroup"
)

type App struct {
	config *config.Config
	db     *dbpg.DB
	logger logger.Logger
	router router.Router
	server *http.Server
}

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

	b.router = router.New(b.logger, b.config.ServerConfig.Mode)

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

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("starting app...")

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		a.logger.Info("starting server...")
		if err := a.server.ListenAndServe(); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		a.logger.Info("received signal... shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			config.DefaultShutdownTimeout)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.logger.Info("server shutdown complete")
		return nil
	})

	//a.logger.Info(fmt.Sprintf("cfg = %v", a.config))

	return g.Wait()
}

func (a *App) Close() error {
	var errs []error

	if a.db != nil {
		if err := a.db.Master.Close(); err != nil {
			errs = append(errs, fmt.Errorf("db close: %w", err))
		}
	}

	return errors.Join(errs...)
}
