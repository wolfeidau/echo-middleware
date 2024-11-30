package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// ZeroLogRequestLogConfig defines the config for the request logger middleware.
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

			start := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}
			r := c.Request()
			w := c.Response()

			log.Ctx(r.Context()).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", w.Status).
				Int64("size", w.Size).
				Int64("dur", time.Since(start).Milliseconds()).
				Str("ua", r.UserAgent()).
				Str("remote-addr", r.RemoteAddr).
				Str("referer", r.Referer()).
				Msg("request")

			return err
		}
	}
}
