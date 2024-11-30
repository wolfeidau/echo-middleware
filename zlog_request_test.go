package middleware

import (
	"bytes"
	"encoding/json"
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

	req := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := ZeroLogWithConfig(ZeroLogConfig{Output: buf, Caller: true})(ZeroLogRequestLog()(handler))

	err := h(c)
	assert.NoError(err)

	var m map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &m)
	assert.NoError(err)

	// assert the map has the expected keys
	assert.Contains(m, "method")
	assert.Contains(m, "url")
	assert.Contains(m, "status")
	assert.Contains(m, "size")
	assert.Contains(m, "dur")
	assert.Contains(m, "ua")
	assert.Contains(m, "remote-addr")
	assert.Contains(m, "referer")
}
