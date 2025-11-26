package middleware

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/url-shortener/internal/logger"
)

func RequestLogger(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		msg := fmt.Sprintf(
			"%s | %3d | %13v | %15s | %-7s %s",
			start.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			path,
		)

		l.Info(msg)
	}
}

func Recovery(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "panic: %v\n", rec)

				// stacktrace
				stack := make([]byte, 4<<10)
				n := runtime.Stack(stack, false)
				buf.Write(stack[:n])

				l.Error(buf.String())
				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}
