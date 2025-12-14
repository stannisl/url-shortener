package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

	// TODO handle fix

	res, err := ah.analyticsService.GetStats(c.Request.Context(), shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dto.AnalyticsSummary{}.FromModel(res)})
}
