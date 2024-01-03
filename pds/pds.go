package pds

import (
	"context"
	"net/http"
	"time"

	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxlogger/logger"
	"github.com/ecumenos/fxecumenos/fxpostgres/postgres"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ServiceName    fxecumenos.ServiceName = "pds"
	ServiceVersion fxecumenos.Version     = "v0.0.0"
)

type Config struct {
	Addr        string
	Prod        bool
	PostgresURL string
}

type PDS struct {
	server          *http.Server
	logger          *zap.Logger
	responseFactory fxrf.Factory
	pg              *postgres.Driver
}

func New(cfg *Config) (*PDS, error) {
	l, err := newLogger(cfg.Prod)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	responseFactory := fxrf.NewFactory(l, &fxrf.Config{WriteLogs: !cfg.Prod}, ServiceVersion)
	pds := &PDS{
		logger:          l,
		responseFactory: responseFactory,
		pg:              driver,
	}

	router := mux.NewRouter()
	enrichContext := NewEnrichContextMiddleware(l, responseFactory)
	recovery := NewRecoverMiddleware(l, responseFactory)
	router.Use(mux.MiddlewareFunc(enrichContext))
	router.HandleFunc("/api/info", pds.Info).Methods(http.MethodGet)
	router.HandleFunc("/api/health", pds.Health).Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(mux.MiddlewareFunc(recovery))

	pds.server = &http.Server{
		Addr:         cfg.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      http.TimeoutHandler(router, 30*time.Second, "something went wrong"),
	}

	return pds, nil
}

func newLogger(prod bool) (*zap.Logger, error) {
	var l *zap.Logger
	var err error
	if prod {
		l, err = logger.NewProductionLogger(string(ServiceName))
	} else {
		l, err = logger.NewDevelopmentLogger(string(ServiceName))
	}
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(l)

	return l, nil
}

func (pds *PDS) Start(ctx context.Context) error {
	if err := pds.pg.Ping(ctx); err != nil {
		return err
	}
	pds.logger.Info("postgres is started")

	return pds.server.ListenAndServe()
}

func (pds *PDS) Shutdown(ctx context.Context) error {
	_ = pds.logger.Sync()
	if err := pds.server.Shutdown(ctx); err != nil {
		return err
	}
	pds.logger.Info("http server was shutted down")

	pds.pg.Close()
	pds.logger.Info("postgres connections was closed")

	return nil
}
