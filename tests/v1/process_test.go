package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"random-exporters/pkg/middleware"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		method       string
		expectedCode int
	}{
		{
			name:         "List processes",
			path:         "/v1/openmetrics/processes",
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
			startTime := time.Now()
			router.ServeHTTP(w, req)
			endTime := time.Now()
			code := w.Code

			if test.expectedCode != code {
				t.Errorf("Expected code %d got %d", test.expectedCode, code)
			}
			//t.Log(w.Body.String())
			t.Logf("Time taken: %d ms", endTime.Sub(startTime).Milliseconds())
		})
	}
}

func TestTopNProcesses(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		method             string
		expectedCode       int
		expectedNumResults int
	}{
		{
			name:               "List processes",
			path:               "/v1/prom/processes?top=20",
			method:             "GET",
			expectedCode:       http.StatusOK,
			expectedNumResults: 20,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			e := gin.New()
			router := middleware.GenerateRouter(e)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, test.path, nil)
			startTime := time.Now()
			router.ServeHTTP(w, req)
			endTime := time.Now()
			code := w.Code

			if test.expectedCode != code {
				t.Errorf("Expected code %d got %d", test.expectedCode, code)
			}
			//t.Log(w.Body.String())
			t.Logf("Time taken: %d ms", endTime.Sub(startTime).Milliseconds())

			body := w.Body.String()
			totalLines := strings.Count(body, "\n")

			if totalLines-7 != test.expectedNumResults*2 {
				t.Errorf("Expected %d results got %d", test.expectedNumResults*2, totalLines-7)
			}
		})
	}
}
