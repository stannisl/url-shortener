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

// Redirect перенаправляет на оригинальный URL
// @Summary      Редирект по короткой ссылке
// @Description  Перенаправляет пользователя на оригинальный URL и записывает аналитику
// @Tags         URLs
// @Produce      html
// @Param        short_url  path  string  true  "Короткий код ссылки"  example(abc123XYZ)
// @Success      302  {string}  string  "Редирект на оригинальный URL"
// @Failure      404  {object}  dto.ErrorResponse  "Ссылка не найдена"
// @Failure      500  {object}  dto.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /s/{short_url} [get]
func (h *UrlHandler) Redirect(c *ginext.Context) {
	shortUrl := c.Param("short_url")
	shortUrl = strings.Trim(shortUrl, "\\/\"") // Триммим \, /, "

	if shortUrl == "" {
		//SetError(c, domain.NewBadRequestError( ,domain.Err)) ??
		return
	}

	urlModel, err := h.urlService.GetOriginUrl(c.Request.Context(), shortUrl)
	if err != nil {
		SetError(c, HandleServiceError(err, "Short URL not found"))
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

// CreateShortUrl создаёт короткую ссылку
// @Summary      Создать короткую ссылку
// @Description  Создаёт сокращённый URL для переданного оригинального адреса
// @Tags         URLs
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUrlRequest  true  "Данные для создания ссылки"
// @Success      201      {object}  dto.APIResponse{data=dto.ShortenURLResponse}  "Ссылка успешно создана"
// @Failure      400      {object}  dto.ErrorResponse  "Неверный формат запроса"
// @Failure      409      {object}  dto.ErrorResponse  "Конфликт (ссылка уже существует)"
// @Failure      422      {object}  dto.ErrorResponse  "Ошибка валидации"
// @Failure      500      {object}  dto.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /shorten [post]
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
