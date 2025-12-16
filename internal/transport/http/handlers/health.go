package handlers

import (
	"github.com/wb-go/wbf/ginext"
)

type Handler struct{}

func NewHealthHandler() *Handler {
	return &Handler{}
}

// Health проверка здоровья сервиса
// @Summary      Health check
// @Description  Проверяет доступность сервиса
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string  "Сервис работает"
// @Router       /health [get]
func (h *Handler) Health(c *ginext.Context) {
	c.JSON(200, ginext.H{"status": "ok"})
}
