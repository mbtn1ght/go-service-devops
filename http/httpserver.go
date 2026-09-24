package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
}

func New(port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Hello)
	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/favicon.ico", Favicon)

	s := &Server{
		server: &http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: mux,
		},
	}

	go s.start()
	log.Info().Msgf("http server listening on port %s", port)
	return s

}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("http server: Shutdown")
	}

	log.Info().Msg("http server: closed")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("request received")
	w.Write([]byte("hello"))
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server: ListenAndServe")
	}
}
