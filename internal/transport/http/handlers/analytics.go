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

// GetAnalytics возвращает статистику по ссылке
// @Summary      Получить аналитику
// @Description  Возвращает статистику переходов по короткой ссылке
// @Tags         Analytics
// @Produce      json
// @Param        short_url  path      string  true  "Короткий код ссылки"  example(abc123XYZ)
// @Success      200        {object}  dto.APIResponse{data=dto.AnalyticsSummaryResponse}  "Статистика получена"
// @Failure      404        {object}  dto.ErrorResponse  "Ссылка не найдена"
// @Failure      500        {object}  dto.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /analytics/{short_url} [get]
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
