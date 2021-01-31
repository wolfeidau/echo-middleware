package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// ZeroLogRequestLogConfig defines the config for the request logger middleware
type ZeroLogRequestLogConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// ZeroLogRequestLog returns a request logger middleware with default config.
func ZeroLogRequestLog() echo.MiddlewareFunc {
	return ZeroLogRequestLogWithConfig(ZeroLogRequestLogConfig{})
}

// ZeroLogRequestLogWithConfig returns a request logger middleware with config.
func ZeroLogRequestLogWithConfig(config ZeroLogRequestLogConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}

			log.Ctx(c.Request().Context()).Info().Fields(map[string]interface{}{
				"path":   req.URL.Path,
				"method": req.Method,
				"dur":    time.Since(start).String(),
				"status": res.Status,
				"length": res.Size,
				"ip":     c.RealIP(),
			}).Msg("processed request")

			return err
		}
	}
}
