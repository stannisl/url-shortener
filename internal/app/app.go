package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/dbpg"
	"golang.org/x/sync/errgroup"

	"github.com/stannisl/url-shortener/internal/config"
	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/transport/http/router"
)

type App struct {
	config *config.Config
	db     *dbpg.DB
	logger logger.Logger
	router router.Router
	server *http.Server
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("starting app...")

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		a.logger.Info("starting server...")
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
