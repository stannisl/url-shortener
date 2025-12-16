package dto

import "github.com/stannisl/url-shortener/internal/domain"

type AnalyticsResponse struct {
	TotalClicks    int `json:"total_clicks"`
	UniqueVisitors int `json:"unique_visitors"`
}

func (ar AnalyticsResponse) FromModel(model *domain.AnalyticsSummary) *AnalyticsResponse {
	return &AnalyticsResponse{
		TotalClicks:    model.TotalClicks,
		UniqueVisitors: model.UniqueVisitors,
	}
}

// ShortenURLResponse - ответ с сокращённым URL
// @Description Ответ с информацией о сокращённой ссылке
type ShortenURLResponse struct {
	OriginalURL string `json:"origin_url"  example:"https://www.google.com"`
	ShortenUrl  string `json:"shorten_url" example:"abc123XYZ"`
}

func (ar ShortenURLResponse) FromModel(model *domain.ShortenUrlModel) ShortenURLResponse {
	return ShortenURLResponse{
		OriginalURL: model.OriginalUrl,
		ShortenUrl:  model.ShortCode,
	}
}

// AnalyticsSummaryResponse - статистика по ссылке
// @Description Аналитика переходов по сокращённой ссылке
type AnalyticsSummaryResponse struct {
	TotalClicks    int `json:"total_clicks" example:"150"`
	UniqueVisitors int `json:"unique_visitors" example:"89"`
}

func (ar AnalyticsSummaryResponse) FromModel(model *domain.AnalyticsSummary) AnalyticsSummaryResponse {
	return AnalyticsSummaryResponse{
		TotalClicks:    model.TotalClicks,
		UniqueVisitors: model.UniqueVisitors,
	}
}

// ErrorResponse - ответ при ошибке
// @Description Ответ при возникновении ошибки
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Resource not found"`
	Code    int    `json:"code,omitempty" example:"404"`
}

type APIResponse struct {
	Data any `json:"data"`
}
