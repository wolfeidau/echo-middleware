package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func handler(c echo.Context) error {
	return c.NoContent(200)
}

func TestZeroLogRequestWithConfig(t *testing.T) {
	assert := require.New(t)

	buf := new(bytes.Buffer)

	e := echo.New()

	e.Use()

	req := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ZeroLogWithConfig(ZeroLogConfig{Output: buf, Caller: true})(ZeroLogRequestLog()(handler))

	err := h(c)
	assert.NoError(err)
	assert.Contains(buf.String(), "ip")
	assert.Contains(buf.String(), "method")
	assert.Contains(buf.String(), "status")
}
