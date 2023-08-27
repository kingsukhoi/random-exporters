package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/showwin/speedtest-go/speedtest"
	"net/http"
	"random-exporters/pkg/formatter/OpenMetrics"
	"strconv"
	"strings"
)

var speedtestClient = speedtest.New()

func GetSpeedTestServersHandler(c *gin.Context) {
	servers, _ := speedtestClient.FetchServers()
	c.JSON(http.StatusOK, servers)
}

func SpeedTestHandler(c *gin.Context) {
	serversString, _ := c.GetQuery("servers")

	servers := strings.Split(serversString, ",")

	if servers[0] == "" {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("servers query is required"))
		return
	}

	sb := strings.Builder{}

	for _, s := range servers {
		server, _ := speedtestClient.FetchServerByID(s)

		err := server.PingTestContext(c.Request.Context(), nil)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		err = server.DownloadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		err = server.UploadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		upload := OpenMetrics.Metric{
			Type:   OpenMetrics.Gauge,
			Name:   "speedtest_upload_mbps",
			Labels: map[string]string{"server": s},
			Value:  strconv.FormatFloat(server.ULSpeed, 'f', -1, 64),
			Help:   "Upload speed for a server",
		}
		download := OpenMetrics.Metric{
			Type:   OpenMetrics.Gauge,
			Name:   "speedtest_download_mbps",
			Labels: map[string]string{"server": s},
			Value:  strconv.FormatFloat(server.DLSpeed, 'f', -1, 64),
			Help:   "Download speed for a server",
		}
		latency := OpenMetrics.Metric{
			Type:   OpenMetrics.Gauge,
			Name:   "speedtest_latency_ms",
			Labels: map[string]string{"server": s},
			Value:  strconv.FormatInt(server.Latency.Milliseconds(), 10),
			Help:   "Latency",
		}

		sb.WriteString(upload.String() + "\n")
		sb.WriteString(download.String() + "\n")
		sb.WriteString(latency.String() + "\n")

	}

	c.Data(http.StatusOK, "application/openmetrics-text; version=1.0.0; charset=utf-8", []byte(sb.String()))
}
