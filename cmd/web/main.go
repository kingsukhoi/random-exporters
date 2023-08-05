package main

import (
	"custom-exporters/pkg/formatter"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/showwin/speedtest-go/speedtest"
	"log"
	"net/http"
	"strconv"
)

var speedtestClient = speedtest.New()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"serving": true})
	})

	r.GET("/servers", func(c *gin.Context) {
		servers, _ := speedtestClient.FetchServers()
		c.JSON(http.StatusOK, servers)
	})

	r.GET("/speedtest", func(c *gin.Context) {
		s, _ := speedtestClient.FetchServerByID("4392")
		log.Println(speedtestClient.FetchServers())

		err := s.PingTestContext(c.Request.Context(), nil)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		err = s.DownloadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		err = s.UploadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		upload := formatter.Gauge{
			Name:   "speedtest_upload_mbps",
			Labels: map[string]string{"server": "4392"},
			Value:  strconv.FormatFloat(s.ULSpeed, 'f', -1, 64),
			Help:   "Upload speed for a server",
		}
		download := formatter.Gauge{
			Name:   "speedtest_download_mbps",
			Labels: map[string]string{"server": "4392"},
			Value:  strconv.FormatFloat(s.DLSpeed, 'f', -1, 64),
			Help:   "Download speed for a server",
		}
		latency := formatter.Gauge{
			Name:   "speedtest_latency_ms",
			Labels: map[string]string{"server": "4392"},
			Value:  strconv.FormatInt(s.Latency.Milliseconds(), 10),
			Help:   "Latency",
		}

		ul, _ := upload.PrintPrometheusGauge()
		dl, _ := download.PrintPrometheusGauge()
		l, _ := latency.PrintPrometheusGauge()

		c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(fmt.Sprintf("%s\n%s\n%s", ul, dl, l)))

	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
