package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"random-exporters/pkg/middleware"
	"testing"
)

func TestWeather(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		method       string
		expectedCode int
	}{
		{
			name:         "Get Weather",
			path:         "/v1/openmetrics/weather?loc=m4g",
			method:       "GET",
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			e := gin.New()
			router := middleware.GenerateRouter(e)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, test.path, nil)
			router.ServeHTTP(w, req)
			code := w.Code

			if test.expectedCode != code {
				t.Errorf("Expected code %d got %d", test.expectedCode, code)
			}
			t.Log(w.Body.String())
		})
	}
}
