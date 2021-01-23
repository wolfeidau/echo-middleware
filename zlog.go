package middleware

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ZeroLogConfig struct {
	Level  zerolog.Level
	Output io.Writer
	Fields map[string]interface{}
}

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

			zl := zerolog.New(cfg.Output).Level(cfg.Level).With().Fields(cfg.Fields).Caller().Stack().Logger()

			c.SetRequest(c.Request().WithContext(zl.WithContext(ctx)))

			return next(c)
		}

	}
}
