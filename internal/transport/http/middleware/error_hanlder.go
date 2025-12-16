package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/url-shortener/internal/domain"
	"github.com/stannisl/url-shortener/internal/logger"
	"github.com/stannisl/url-shortener/internal/transport/http/dto"
	"github.com/wb-go/wbf/ginext"
)

func ErrorHandler(l logger.Logger) ginext.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, l, err)
			return
		}
	}
}

func handleError(c *ginext.Context, l logger.Logger, err error) {
	var appError *domain.Error

	if errors.As(err, &appError) {
		if appError.Code >= 500 {
			l.ErrorErr(appError.Err, appError.Message)
		}
		c.JSON(appError.Code,
			dto.ErrorResponse{
				Success: false,
				Error:   appError.Message,
				Code:    appError.Code,
			})
		return
	}

	switch {
	case errors.Is(err, domain.ErrNotUnique):
		c.JSON(http.StatusConflict,
			dto.ErrorResponse{
				Success: false,
				Error:   "Resource already exists",
				Code:    http.StatusConflict,
			})
	case errors.Is(err, domain.ErrModelNotFound):
		c.JSON(http.StatusNotFound,
			dto.ErrorResponse{
				Success: false,
				Error:   "Resource not found",
				Code:    http.StatusNotFound,
			})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest,
			dto.ErrorResponse{
				Success: false,
				Error:   "Bad Request",
				Code:    http.StatusBadRequest,
			})
	default:
		l.ErrorErr(appError.Err, appError.Message)
		c.JSON(http.StatusInternalServerError,
			dto.ErrorResponse{
				Success: false,
				Error:   "Internal server error",
				Code:    http.StatusInternalServerError,
			})
	}
}
