package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	gen "github.com/ecumenos/ecumenos/internal/generated/pds"
	"github.com/ecumenos/ecumenos/internal/httputils"
	"github.com/ecumenos/ecumenos/pds/config"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxresponsefactory.Factory
}

type serverParams struct {
	fx.In
	Config    *config.Config
	Logger    *zap.Logger
	ServerInt gen.ServerInterface
}

func NewServer(params serverParams) *Server {
	responseFactory := fxresponsefactory.NewFactory(params.Logger, &fxresponsefactory.Config{WriteLogs: !params.Config.Prod}, config.ServiceVersion)
	s := &Server{
		logger:          params.Logger,
		responseFactory: responseFactory,
	}

	router := mux.NewRouter()
	enrichContext := httputils.NewEnrichContextMiddleware(params.Logger, responseFactory)
	recovery := httputils.NewRecoverMiddleware(params.Logger, responseFactory)
	router.Use(mux.MiddlewareFunc(enrichContext))
	router = gen.HandlerWithOptions(params.ServerInt, gen.GorillaServerOptions{
		BaseRouter:       router,
		ErrorHandlerFunc: httputils.DefaultErrorHandlerFactory(responseFactory),
	}).(*mux.Router)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(mux.MiddlewareFunc(recovery))
	s.server = &http.Server{
		Addr:         params.Config.AppAddr,
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
