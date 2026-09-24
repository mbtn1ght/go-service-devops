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
	port := "8080"
	httpServer := http.New(port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Info().Msg("Shutting down...")
	httpServer.Close()
}
