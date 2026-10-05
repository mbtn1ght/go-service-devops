package main

import (
	"os"
	"os/signal"
	"service_devops/http"
	"service_devops/logger"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	webDir := os.Getenv("WEB_DIR")
	if webDir == "" {
		webDir = "web/dist"
	}
	httpServer := http.New(port, webDir)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Info().Msg("Shutting down...")
	httpServer.Close()
}
