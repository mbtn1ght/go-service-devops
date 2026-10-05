package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
}

func New(port, webDir string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/", SPA(webDir))

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

func SPA(webDir string) http.HandlerFunc {
	files := http.FileServer(http.Dir(webDir))

	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("request received")

		path := filepath.Join(webDir, filepath.Clean(strings.TrimPrefix(r.URL.Path, "/")))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server: ListenAndServe")
	}
}
