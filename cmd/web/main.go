package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log"
	"log/slog"
	"os"
	"random-exporters/pkg/middleware"
)

func main() {
	r := gin.New()

	if gin.Mode() == gin.ReleaseMode {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		r.Use(sloggin.New(logger))
	} else {
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())
	middleware.GenerateRouter(r)
	pprof.Register(r)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
