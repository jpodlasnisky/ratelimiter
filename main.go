// main.go
package main

import (
	"log"
	"net/http"

	"github.com/jpodlasnisky/ratelimiter/config"
	"github.com/jpodlasnisky/ratelimiter/infra/web/middleware"
	"github.com/jpodlasnisky/ratelimiter/infra/web/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	rateLimiter := server.SetupRateLimiter(cfg)

	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	rateLimitMiddleware := middleware.RateLimitMiddleware(mux, rateLimiter)

	loggingMiddleware := middleware.LoggingMiddleware(rateLimitMiddleware)

	srv := server.New(cfg.WebPort, loggingMiddleware)

	go func() {
		log.Println("Servidor HTTP iniciado na porta:", cfg.WebPort)
		srv.Start()
	}()

	server.WaitForShutdown(srv.Server)
}
