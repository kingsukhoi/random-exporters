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

	v1Prom := v1.Group("/prom")
	v1Prom.GET("/speedtest", routes.SpeedTestHandler)

	v1Prom.GET("/processes", routes.ProcessHandler)

	return r
}
