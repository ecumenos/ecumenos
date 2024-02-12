package admin

import (
	"fmt"
	"net/http"

	"github.com/ecumenos/ecumenos/internal/docs"
	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	gen "github.com/ecumenos/ecumenos/internal/generated/orbissociusadmin"
	"github.com/ecumenos/ecumenos/internal/openapi"
	"github.com/ecumenos/ecumenos/orbissocius/config"
	"github.com/ecumenos/ecumenos/orbissocius/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	responseFactory fxresponsefactory.Factory
	Service         *service.Service
	selfURL         string
}

type handlerParams struct {
	fx.In
	Service *service.Service
	Config  *config.Config
	Logger  *zap.Logger
}

func NewHandler(params handlerParams) gen.ServerInterface {
	responseFactory := fxresponsefactory.NewFactory(params.Logger, &fxresponsefactory.Config{WriteLogs: !params.Config.Prod}, config.ServiceVersion)

	return &handler{
		responseFactory: responseFactory,
		Service:         params.Service,
		selfURL:         params.Config.AdminSelfURL,
	}
}

func (h *handler) GetDocs(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	_, _ = rw.Write(docs.OrbisSociusAdminDocs(h.selfURL))
}

func (h *handler) GetHealth(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, gen.GetHealthData{Ok: true}) //nolint:errcheck
}

func (h *handler) GetInfo(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result := h.Service.PingServices(ctx)
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, gen.GetInfoData{ //nolint:errcheck
		Name:    string(config.ServiceName),
		Version: string(config.ServiceVersion),
		Deps:    &map[string]interface{}{"postgres": result.PostgresIsRunning},
	})
}

func (h *handler) GetSpecs(rw http.ResponseWriter, r *http.Request) {
	filename := "orbis_socius_admin_specs.yaml"
	rw.Header().Add("Content-Type", "application/openapi+yaml")
	rw.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	_, _ = rw.Write(openapi.OrbisSociusAdminSpec(h.selfURL))
}
