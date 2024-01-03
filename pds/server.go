package pds

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	pds             *PDS
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxrf.Factory
}

func NewServer(cfg *Config, pds *PDS, l *zap.Logger) *Server {
	responseFactory := fxrf.NewFactory(l, &fxrf.Config{WriteLogs: !cfg.Prod}, ServiceVersion)
	s := &Server{
		pds:             pds,
		logger:          l,
		responseFactory: responseFactory,
	}

	router := mux.NewRouter()
	enrichContext := NewEnrichContextMiddleware(l, responseFactory)
	recovery := NewRecoverMiddleware(l, responseFactory)
	router.Use(mux.MiddlewareFunc(enrichContext))
	router.HandleFunc("/api/info", s.Info).Methods(http.MethodGet)
	router.HandleFunc("/api/health", s.Health).Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(mux.MiddlewareFunc(recovery))
	s.server = &http.Server{
		Addr:         cfg.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      http.TimeoutHandler(router, 30*time.Second, "something went wrong"),
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("http server was shutted down")

	return nil
}
