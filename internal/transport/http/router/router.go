package router

import (
	"net/http"

	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/service"
	"github.com/stannisl/url-shortener/internal/transport/http/handlers"
	"github.com/stannisl/url-shortener/internal/transport/http/middleware"
	"github.com/wb-go/wbf/ginext"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	http.Handler
}

type ginRouter struct {
	router *ginext.Engine
}

func New(
	l logger.Logger,
	urlService *service.UrlService,
	analyticsService *service.AnalyticsService,
	ginMode string,
) Router {
	router := ginext.New(ginMode)
	registerMiddlewares(l, router)

	healthHandler := handlers.NewHealthHandler()
	urlHandler := handlers.NewUrlHandler(urlService, analyticsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	health := router.Group("/health")
	{
		health.GET("", healthHandler.Health)
	}
	
	api := router.Group("")
	{
		api.GET("/analytics/:short_url", analyticsHandler.GetAnalytics)
		api.POST("/shorten", urlHandler.CreateShortUrl)
		api.GET("/s/:short_url", urlHandler.Redirect)
	}

	return &ginRouter{
		router: router,
	}
}

func registerMiddlewares(l logger.Logger, engine *ginext.Engine) {
	engine.Use(middleware.RequestLogger(l))
	engine.Use(middleware.Recovery(l))
	engine.Use(middleware.ErrorHandler(l))
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
