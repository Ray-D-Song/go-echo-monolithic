package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// Logger returns the logger middleware
func Logger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogMethod:    true,
		LogLatency:   true,
		LogRequestID: true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("request",
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency", v.Latency),
					zap.String("request_id", v.RequestID),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("user_agent", v.UserAgent),
				)
			} else {
				logger.Error("request",
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Duration("latency", v.Latency),
					zap.String("request_id", v.RequestID),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("user_agent", v.UserAgent),
					zap.Error(v.Error),
				)
			}
			return nil
		},
	})
}

// AccessLogger returns a simpler access logger middleware for development
func AccessLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human} ${remote_ip} ${user_agent}\n",
		CustomTimeFormat: time.RFC3339,
	})
}