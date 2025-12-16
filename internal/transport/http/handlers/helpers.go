package handlers

import (
	"errors"

	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/wb-go/wbf/ginext"
)

func SetError(c *ginext.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func HandleServiceError(err error, notFoundMsg string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrModelNotFound):
		return domain.NewNotFoundError(notFoundMsg)
	case errors.Is(err, domain.ErrNotUnique):
		return domain.NewConflictError("Resource already exists")
	default:
		return domain.NewInternalError("Internal server error", err)
	}
}
