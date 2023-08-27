package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"random-exporters/pkg/middleware"
	"strings"
	"testing"
)

func TestGetServerList(t *testing.T) {
	e := gin.New()
	router := middleware.GenerateRouter(e)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/servers", nil)
	router.ServeHTTP(w, req)
	serverList := w.Body.String()

	if strings.TrimSpace(serverList) == "" {
		t.Error("server list is empty")
		t.Fail()
	}
}

func TestSpeedTest(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		method       string
		expectedCode int
	}{
		{
			name:         "Do Speed Test",
			path:         "/v1/prom/speedtest?servers=4392",
			method:       "GET",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty server list",
			path:         "/v1/prom/speedtest",
			method:       "GET",
			expectedCode: http.StatusBadRequest,
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
