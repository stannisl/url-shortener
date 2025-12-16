package main

import (
	"context"

	"github.com/stannisl/url-shortener/internal/app"
	"github.com/stannisl/url-shortener/internal/config"
	"github.com/stannisl/url-shortener/internal/infra/postgres"
	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/wb-go/wbf/zlog"

	_ "github.com/stannisl/url-shortener/docs" // Импорт сгенерированных docs
)

// @title           URL Shortener API
// @version         1.0
// @description     Сервис для сокращения URL и сбора аналитики переходов

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @schemes http https
func main() {
	// initializing logger
	zlog.InitConsole()
	adaptedLogger := logger.NewZLogger(zlog.Logger)

	// getting application config
	cfg, err := config.LoadConfig("./config.yaml", "./.env")
	if err != nil {
		adaptedLogger.FatalErr(err, "failed to load config")
	}

	// getting postgres db with loaded configuration
	db, err := postgres.NewDB(cfg.PostgresConfig)
	if err != nil {
		adaptedLogger.FatalErr(err, "failed to connect to database")
	}

	// building application with necessary dependencies
	application, err := app.NewBuilder().
		WithConfig(cfg).
		WithPostgresDB(db).
		WithLogger(adaptedLogger).
		Build()
	if err != nil {
		adaptedLogger.FatalErr(err, "failed to build app")
	}

	if err := application.Run(context.Background()); err != nil {
		adaptedLogger.FatalErr(err, "application error")
	}

	defer application.Close()
}
