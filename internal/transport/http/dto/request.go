package dto

// CreateUrlRequest - запрос на создание короткой ссылки
// @Description Запрос на создание сокращённого URL
type CreateUrlRequest struct {
	OriginUrl string `json:"origin_url" binding:"required,url" example:"https://www.google.com"`
}
