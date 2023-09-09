package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Unix epoch time.
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

// NoCacheConfig used to configure the no cache middleware.
type NoCacheConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// NoCache returns a middleware which sets the no cache headers with default configuration.
func NoCache() echo.MiddlewareFunc {
	return NoCacheWithConfig(NoCacheConfig{})
}

// NoCacheWithConfig returns a middleware which sets a number of http headers to ensure the resource is not cached.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, private, max-age=0
//	X-Accel-Expires: 0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCacheWithConfig(config NoCacheConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			for k, v := range noCacheHeaders {
				c.Response().Header().Set(k, v)
			}

			return next(c)
		}
	}
}
