package middleware

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// ZeroLogConfig used to configure the zerolog echo middleware
type ZeroLogConfig struct {
	Caller bool
	Level  zerolog.Level
	Output io.Writer
	Fields map[string]interface{}
}

// ZeroLogWithConfig setup and return an echo middleware with zerolog logger available from the context.Context.
func ZeroLogWithConfig(cfg ZeroLogConfig) echo.MiddlewareFunc {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	if cfg.Fields == nil {
		cfg.Fields = make(map[string]interface{})
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			zc := zerolog.New(cfg.Output).Level(cfg.Level).With().Fields(cfg.Fields)

			if cfg.Caller {
				zc = zc.Caller()
			}

			zl := zc.Logger()

			c.SetRequest(c.Request().WithContext(zl.WithContext(ctx)))

			return next(c)
		}

	}
}
