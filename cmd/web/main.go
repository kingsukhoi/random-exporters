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
	//todo add setting to switch between text and json (for server)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())
	middleware.GenerateRouter(r)
	pprof.Register(r)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
