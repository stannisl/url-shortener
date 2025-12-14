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

type AnalyticsSummary struct {
	TotalClicks    int `json:"total_clicks"`
	UniqueVisitors int `json:"unique_visitors"`
}

func (ar AnalyticsSummary) FromModel(model *domain.AnalyticsSummary) AnalyticsSummary {
	return AnalyticsSummary{
		TotalClicks:    model.TotalClicks,
		UniqueVisitors: model.UniqueVisitors,
	}
}
