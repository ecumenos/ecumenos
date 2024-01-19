package orbissocius

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AdminServer struct {
	o               *OrbisSocius
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxresponsefactory.Factory
}

func NewAdminServer(cfg *Config, o *OrbisSocius, l *zap.Logger) *AdminServer {
	responseFactory := fxresponsefactory.NewFactory(l, &fxresponsefactory.Config{WriteLogs: !cfg.Prod}, ServiceVersion)
	s := &AdminServer{
		o:               o,
		logger:          l,
		responseFactory: responseFactory,
	}

	router := mux.NewRouter()
	enrichContext := NewEnrichContextMiddleware(l, responseFactory)
	recovery := NewRecoverMiddleware(l, responseFactory)
	router.Use(mux.MiddlewareFunc(enrichContext))
	router.HandleFunc("/api/admin/info", s.Info).Methods(http.MethodGet)
	router.HandleFunc("/api/admin/health", s.Health).Methods(http.MethodGet)
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
