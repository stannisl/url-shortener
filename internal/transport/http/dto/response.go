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

type ShortenURLResponse struct {
	OriginalURL string `json:"origin_url"`
	ShortenUrl  string `json:"shorten_url"`
}

func (ar ShortenURLResponse) FromModel(model *domain.ShortenUrlModel) ShortenURLResponse {
	return ShortenURLResponse{
		OriginalURL: model.OriginalUrl,
		ShortenUrl:  model.ShortCode,
	}
}

type AnalyticsSummaryResponse struct {
	TotalClicks    int `json:"total_clicks"`
	UniqueVisitors int `json:"unique_visitors"`
}

func (ar AnalyticsSummaryResponse) FromModel(model *domain.AnalyticsSummary) AnalyticsSummaryResponse {
	return AnalyticsSummaryResponse{
		TotalClicks:    model.TotalClicks,
		UniqueVisitors: model.UniqueVisitors,
	}
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}
