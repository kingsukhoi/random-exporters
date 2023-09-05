package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/showwin/speedtest-go/speedtest"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"random-exporters/pkg/middleware"
	"strings"
	"testing"
	"time"
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

func TestAgainstOfficialTool(t *testing.T) {
	out, err := exec.Command("speedtest-go", "-s", "4392", "--json", "--force-http-ping").Output()
	if err != nil {
		t.Fatal(err)
		return
	}
	officialResult := &OfficialSpeedTest{}

	json.Unmarshal(out, officialResult)
	officialDL := officialResult.Servers[0].DlSpeed
	officialUl := officialResult.Servers[0].UlSpeed

	t.Logf("Ping %d", officialResult.Servers[0].TestDuration.Ping)
	t.Logf("Download %f", officialDL)
	t.Logf("Upload %f", officialUl)

	var speedtestClient = speedtest.New()
	server, err := speedtestClient.FetchServerByID("4392")
	err = server.TestAll()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Go Dl %f", server.DLSpeed)
	t.Logf("Go Ul %f", server.ULSpeed)

	if server.DLSpeed < officialDL*0.9 || server.DLSpeed > officialDL*1.10 {
		t.Errorf("(Download) Go program got %f, official got %f. Difference is greater than 10%%", server.DLSpeed, officialDL)
	}

	if server.ULSpeed < officialUl*0.9 || server.ULSpeed > officialUl*1.10 {
		t.Errorf("(Upload) Go program got %f, official got %f. Difference is greater than 10%%", server.DLSpeed, officialDL)
	}
}

type OfficialSpeedTest struct {
	Timestamp string `json:"timestamp"`
	UserInfo  struct {
		IP  string `json:"IP"`
		Lat string `json:"Lat"`
		Lon string `json:"Lon"`
		Isp string `json:"Isp"`
	} `json:"user_info"`
	Servers []struct {
		Url          string  `json:"url"`
		Lat          string  `json:"lat"`
		Lon          string  `json:"lon"`
		Name         string  `json:"name"`
		Country      string  `json:"country"`
		Sponsor      string  `json:"sponsor"`
		Id           string  `json:"id"`
		Url2         string  `json:"url_2"`
		Host         string  `json:"host"`
		Distance     float64 `json:"distance"`
		Latency      int     `json:"latency"`
		MaxLatency   int     `json:"max_latency"`
		MinLatency   int     `json:"min_latency"`
		Jitter       int     `json:"jitter"`
		DlSpeed      float64 `json:"dl_speed"`
		UlSpeed      float64 `json:"ul_speed"`
		TestDuration struct {
			Ping     time.Duration `json:"ping"`
			Download time.Duration `json:"download"`
			Upload   time.Duration `json:"upload"`
			Total    time.Duration `json:"total"`
		} `json:"test_duration"`
	} `json:"servers"`
}
