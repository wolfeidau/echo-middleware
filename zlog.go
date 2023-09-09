package middleware

import (
	"io"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// ZeroLogConfig used to configure the zerolog echo middleware.
type ZeroLogConfig struct {
	Caller bool
	Level  zerolog.Level
	Output io.Writer
	Logger zerolog.Logger
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

			if reflect.ValueOf(cfg.Logger).IsZero() {
				zc := zerolog.New(cfg.Output).Level(cfg.Level).With().Fields(cfg.Fields)
				if cfg.Caller {
					zc = zc.Caller()
				}

				cfg.Logger = zc.Logger()
			}

			c.SetRequest(c.Request().WithContext(cfg.Logger.WithContext(ctx)))

			return next(c)
		}
	}
}
