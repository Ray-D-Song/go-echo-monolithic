package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestID returns the request ID middleware
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			// Use Echo's default UUID generator
			return middleware.DefaultRequestIDConfig.Generator()
		},
	})
}