package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stannisl/url-shortener/internal/service"
	"github.com/stannisl/url-shortener/internal/transport/http/dto"
	"github.com/wb-go/wbf/ginext"
)

type UrlHandler struct {
	urlService       *service.UrlService
	analyticsService *service.AnalyticsService
}

func NewUrlHandler(urlService *service.UrlService, analyticsService *service.AnalyticsService) *UrlHandler {
	return &UrlHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

func (h *UrlHandler) Redirect(c *ginext.Context) {
	shortUrl := c.Param("short_url")
	shortUrl = strings.Trim(shortUrl, "\\/\"") // Триммим \, /, "

	if shortUrl == "" {
		SetError(c, domain.NewBadRequestError("Invalid short_url", domain.ErrModelNotFound))
		return
	}

	urlModel, err := h.urlService.GetOriginUrl(c.Request.Context(), shortUrl)
	if err != nil {
		SetError(c, domain.NewBadRequestError("Short URL not found", err))
		return
	}

	go func() {
		_ = h.analyticsService.HandleRedirect(c.Request.Context(), shortUrl, domain.UrlAnalyticsModel{
			UserAgent:  c.Request.UserAgent(),
			IPAddress:  c.ClientIP(),
			AccessedAt: time.Now(),
		})
	}()

	c.Redirect(http.StatusFound, urlModel.OriginalUrl)
}

func (h *UrlHandler) CreateShortUrl(c *gin.Context) {
	var createUrlRequest dto.CreateUrlRequest

	if err := c.ShouldBindJSON(&createUrlRequest); err != nil {
		SetError(c, domain.NewBadRequestError("Invalid JSON Schema", err))
		return
	}

	if createUrlRequest.OriginUrl == "" {
		SetError(c, domain.NewBadRequestError("Invalid origin_url", domain.ErrInvalidInput))
		return
	}

	shortURL, err := h.urlService.CreateShortenUrl(c.Request.Context(), createUrlRequest.OriginUrl)
	if err != nil {
		SetError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ShortenURLResponse{}.FromModel(shortURL))
}
