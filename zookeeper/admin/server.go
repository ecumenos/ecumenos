package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	"github.com/ecumenos/ecumenos/internal/fxtypes"
	"github.com/ecumenos/ecumenos/internal/httputils"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	service         *service.Service
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxresponsefactory.Factory
	serviceVersion  fxtypes.Version
	serviceName     fxtypes.ServiceName
}

type serverParams struct {
	fx.In

	Config         *config.Config
	Service        *service.Service
	Logger         *zap.Logger
	ServiceVersion fxtypes.Version
	ServiceName    fxtypes.ServiceName
}

func New(p serverParams) *Server {
	responseFactory := fxresponsefactory.NewFactory(p.Logger, &fxresponsefactory.Config{WriteLogs: !p.Config.Prod}, p.ServiceVersion)
	s := &Server{
		service:         p.Service,
		logger:          p.Logger,
		responseFactory: responseFactory,
		serviceVersion:  p.ServiceVersion,
		serviceName:     p.ServiceName,
	}

	r := mux.NewRouter().PathPrefix("/api/admin").Subrouter()
	enrichContext := httputils.NewEnrichContextMiddleware(p.Logger, responseFactory)
	recovery := httputils.NewRecoverMiddleware(p.Logger, responseFactory)
	r.Use(mux.MiddlewareFunc(enrichContext))
	r.HandleFunc("/info", s.Info).Methods(http.MethodGet)
	r.HandleFunc("/health", s.Health).Methods(http.MethodGet)
	r.HandleFunc("/sign-in", s.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/refresh-session", s.RefreshSession).Methods(http.MethodPost)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(mux.MiddlewareFunc(recovery))

	guarded := r.NewRoute().Subrouter()
	auth := httputils.NewAdminAuthorizationMiddleware(p.Logger, responseFactory, p.Service)
	guarded.Use(mux.MiddlewareFunc(auth))
	guarded.HandleFunc("/sign-out", s.SignOut).Methods(http.MethodDelete)

	s.server = &http.Server{
		Addr:         p.Config.AdminAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      http.TimeoutHandler(r, 30*time.Second, "something went wrong"),
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
