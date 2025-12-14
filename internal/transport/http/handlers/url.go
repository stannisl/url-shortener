package handlers

import (
	"log"
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

	urlModel, err := h.urlService.GetOriginUrl(c.Request.Context(), shortUrl)

	// maybe put it in middleware
	if err != nil {
		log.Println(err)
		// handle service errors or parse errors 404 500 and more
		return
	}

	err = h.analyticsService.HandleRedirect(c.Request.Context(), shortUrl, domain.UrlAnalyticsModel{
		UserAgent:  c.Request.UserAgent(),
		IPAddress:  "8.8.8.8",
		AccessedAt: time.Now(),
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(err, urlModel.OriginalUrl)

	c.Redirect(http.StatusFound, urlModel.OriginalUrl)
}

func (h *UrlHandler) CreateShortUrl(c *gin.Context) {
	var createUrlRequest dto.CreateUrlRequest

	if err := c.ShouldBindJSON(&createUrlRequest); err != nil {
		// handle
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.urlService.CreateShortenUrl(c.Request.Context(), createUrlRequest.OriginUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ShortenURLResponse{}.FromModel(shortURL))
}
