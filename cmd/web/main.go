package main

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"random-exporters/pkg/middleware"
	"time"
)

func main() {
	r := gin.New()

	logger, _ := zap.NewProduction()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	middleware.GenerateRouter(r)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
