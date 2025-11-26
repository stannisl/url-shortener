package router

import (
	"net/http"

	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/transport/http/handlers"
	"github.com/stannisl/url-shortener/internal/transport/http/middleware"
	"github.com/wb-go/wbf/ginext"
)

type Router interface {
	http.Handler
}

type ginRouter struct {
	router *ginext.Engine
}

func New(l logger.Logger, ginMode string) Router {
	router := ginext.New(ginMode)
	registerMiddlewares(l, router)

	healthHandler := handlers.NewHealthHandler()

	health := router.Group("/health")
	{
		health.GET("", healthHandler.Health)
	}

	return &ginRouter{
		router: router,
	}
}

func registerMiddlewares(l logger.Logger, engine *ginext.Engine) {
	engine.Use(middleware.RequestLogger(l))
	engine.Use(middleware.Recovery(l))
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
