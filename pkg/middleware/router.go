package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"random-exporters/pkg/routes"
)

func GenerateRouter(r *gin.Engine) *gin.Engine {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"serving": true})
	})

	v1 := r.Group("/v1")
	v1.GET("/servers", routes.GetSpeedTestServersHandler)

	v1om := v1.Group("/openmetrics")
	v1om.GET("/speedtest", routes.SpeedTestHandler)
	v1om.GET("/processes", routes.ProcessHandler)
	v1om.GET("/weather", routes.OpenMetricsWeatherCurrent)

	return r
}
