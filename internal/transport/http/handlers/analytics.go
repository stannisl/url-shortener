package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stannisl/url-shortener/internal/service"
	"github.com/stannisl/url-shortener/internal/transport/http/dto"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: *service,
	}
}

func (ah *AnalyticsHandler) GetAnalytics(c *gin.Context) {
	shortUrl := c.Param("short_url")
	shortUrl = strings.Trim(shortUrl, "/")
	if shortUrl == "" {
		SetError(c, domain.NewBadRequestError("Invalid short_url", domain.ErrInvalidInput))
		return
	}

	res, err := ah.analyticsService.GetStats(c.Request.Context(), shortUrl)
	if err != nil {
		SetError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": dto.AnalyticsSummaryResponse{}.FromModel(res)})
}
