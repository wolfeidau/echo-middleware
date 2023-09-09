package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestZeroLogWithConfig(t *testing.T) {
	assert := require.New(t)

	buf := new(bytes.Buffer)

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ZeroLogWithConfig(ZeroLogConfig{Output: buf, Caller: true})(func(c echo.Context) error {
		log.Ctx(c.Request().Context()).Info().Msg("woops")

		return c.NoContent(200)
	})

	err := h(c)
	assert.NoError(err)
	assert.Contains(buf.String(), "woops")
	assert.Contains(buf.String(), "info")
	assert.Contains(buf.String(), "caller")
}

func TestZeroLogWithConfigLogger(t *testing.T) {
	assert := require.New(t)

	buf := new(bytes.Buffer)

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ZeroLogWithConfig(ZeroLogConfig{Logger: zerolog.New(buf)})(func(c echo.Context) error {
		log.Ctx(c.Request().Context()).Info().Msg("woops")

		return c.NoContent(200)
	})

	err := h(c)
	assert.NoError(err)
	assert.Contains(buf.String(), "woops")
	assert.Contains(buf.String(), "info")
}
