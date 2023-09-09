package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const OriginVerifyHeaderName = "X-Origin-Verify"

// OriginVerifierConfig used to configure the origin authentication middleware.
type OriginVerifierConfig struct {
	// Token used to validate requests coming include the required header
	Token string

	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

// OriginVerifierWithConfig returns a middleware which verifies requests include a `X-Origin-Verify` header
// containing the token configured, requests which fail will be rejected with a 400 bad request status code.
//
// This solution is based on a pattern presented in https://aws.amazon.com/blogs/networking-and-content-delivery/restricting-access-http-api-gateway-lambda-authorizer/
// and uses the same header name.
func OriginVerifierWithConfig(config OriginVerifierConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			headerToken := c.Request().Header.Get(OriginVerifyHeaderName)

			if subtle.ConstantTimeCompare([]byte(config.Token), []byte(headerToken)) != 1 {
				return c.String(http.StatusBadRequest, "Bad Request")
			}

			return next(c)
		}
	}
}
