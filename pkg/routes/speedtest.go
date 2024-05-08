package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/showwin/speedtest-go/speedtest"
	"net/http"
	"random-exporters/pkg/formatter/OpenMetrics"
	"strconv"
	"strings"
)

func GetSpeedTestServersHandler(c *gin.Context) {
	var speedtestClient = speedtest.New()
	servers, _ := speedtestClient.FetchServers()
	c.JSON(http.StatusOK, servers)
}

func SpeedTestHandler(c *gin.Context) {
	var speedtestClient = speedtest.New()
	serversString, _ := c.GetQuery("servers")

	servers := strings.Split(serversString, ",")

	if servers[0] == "" {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("servers query is required"))
		return
	}

	sb := strings.Builder{}

	for _, s := range servers {
		server, err := speedtestClient.FetchServerByID(s)

		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = server.PingTestContext(c.Request.Context(), nil)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		err = server.DownloadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		err = server.UploadTestContext(c.Request.Context())
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		//slog.Info("output from speed test ul", "output", server.DLSpeed)
		//slog.Info("output from speed test ul", "output", server.DLSpeed.Byte(speedtest.UnitTypeBinaryBytes))
		upload := OpenMetrics.Metric{
			Type:   OpenMetrics.Gauge,
			Name:   "speedtest_upload_bps",
			Labels: map[string]string{"server": s},
			Value:  fmt.Sprintf("%f", server.ULSpeed),
			Help:   "Upload speed for a server",
		}
		download := OpenMetrics.Metric{
			Type:   OpenMetrics.Gauge,
			Name:   "speedtest_download_bps",
			Labels: map[string]string{"server": s},
			Value:  fmt.Sprintf("%f", server.DLSpeed),
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
