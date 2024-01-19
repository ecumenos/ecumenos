package zookeeper

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AdminServer struct {
	zookeeper       *Zookeeper
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxresponsefactory.Factory
}

func NewAdminServer(cfg *Config, z *Zookeeper, l *zap.Logger) *AdminServer {
	responseFactory := fxresponsefactory.NewFactory(l, &fxresponsefactory.Config{WriteLogs: !cfg.Prod}, ServiceVersion)
	s := &AdminServer{
		zookeeper:       z,
		logger:          l,
		responseFactory: responseFactory,
	}

	r := mux.NewRouter().PathPrefix("/api/admin").Subrouter()
	enrichContext := NewEnrichContextMiddleware(l, responseFactory)
	recovery := NewRecoverMiddleware(l, responseFactory)
	r.Use(mux.MiddlewareFunc(enrichContext))
	r.HandleFunc("/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/health", s.Health).Methods(http.MethodGet)
	r.HandleFunc("/sign-in", s.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/refresh-session", s.RefreshSession).Methods(http.MethodPost)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(mux.MiddlewareFunc(recovery))

	guarded := r.NewRoute().Subrouter()
	auth := NewAdminAuthorizationMiddleware(l, responseFactory, z)
	guarded.Use(mux.MiddlewareFunc(auth))
	guarded.HandleFunc("/sign-out", s.SignOut).Methods(http.MethodDelete)

	s.server = &http.Server{
		Addr:         cfg.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      http.TimeoutHandler(r, 30*time.Second, "something went wrong"),
	}

	return s
}

func (s *AdminServer) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *AdminServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("http server was shutted down")

	return nil
}
