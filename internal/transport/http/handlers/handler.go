package handlers

import (
	"github.com/wb-go/wbf/ginext"
)

type Handler struct{}

func NewHealthHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(c *ginext.Context) {
	c.JSON(200, ginext.H{"status": "ok"})
}
