package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestNoCacheWithConfig(t *testing.T) {
	assert := require.New(t)

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NoCacheWithConfig(NoCacheConfig{})(func(c echo.Context) error {
		return c.NoContent(200)
	})

	err := h(c)
	assert.NoError(err)
	assert.Equal(200, rec.Result().StatusCode)
	for k, v := range noCacheHeaders {
		assert.Contains(rec.Header().Get(k), v)
	}
}
