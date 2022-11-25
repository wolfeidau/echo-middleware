package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestOriginVerifierWithConfig(t *testing.T) {
	type args struct {
		config OriginVerifierConfig
	}
	tests := []struct {
		name       string
		args       args
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "should accept request with valid verify token",
			args:       args{config: OriginVerifierConfig{Token: "test"}},
			headers:    map[string]string{"X-Origin-Verify": "test"},
			wantStatus: 200,
		},
		{
			name:       "should reject request with no valid verify token",
			args:       args{config: OriginVerifierConfig{Token: "test"}},
			headers:    map[string]string{},
			wantStatus: 400,
		},
		{
			name:       "should reject request with an invalid verify token",
			args:       args{config: OriginVerifierConfig{Token: "test"}},
			headers:    map[string]string{"X-Origin-Verify": "nottest"},
			wantStatus: 400,
		},
		{
			name:       "should skip request which match the skipper",
			args:       args{config: OriginVerifierConfig{Token: "test", Skipper: func(c echo.Context) bool { return true }}},
			headers:    map[string]string{"X-Origin-Verify": "nottest"},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)

			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/login", nil)

			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := OriginVerifierWithConfig(tt.args.config)(func(c echo.Context) error {
				return c.NoContent(200)
			})

			err := h(c)
			assert.NoError(err)
			assert.Equal(tt.wantStatus, rec.Result().StatusCode)
		})
	}
}
