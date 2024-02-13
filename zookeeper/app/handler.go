package app

import (
	"fmt"
	"net/http"

	"github.com/ecumenos/ecumenos/internal/docs"
	f "github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	gen "github.com/ecumenos/ecumenos/internal/generated/zookeeper"
	"github.com/ecumenos/ecumenos/internal/openapi"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	responseFactory f.Factory
	service         *service.Service
	selfURL         string
}

type handlerParams struct {
	fx.In
	Service *service.Service
	Config  *config.Config
	Logger  *zap.Logger
}

func NewHandler(params handlerParams) gen.ServerInterface {
	responseFactory := f.NewFactory(params.Logger, &f.Config{WriteLogs: !params.Config.Prod}, config.ServiceVersion)

	return &handler{
		responseFactory: responseFactory,
		service:         params.Service,
		selfURL:         params.Config.AppSelfURL,
	}
}

func (h *handler) GetDocs(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	_, _ = rw.Write(docs.ZookeeperDocs(h.selfURL))
}

func (h *handler) GetHealth(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, gen.GetHealthData{Ok: true}) //nolint:errcheck
}

func (h *handler) GetInfo(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deps := h.service.PingServices(ctx)
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, gen.GetInfoData{ //nolint:errcheck
		Name:    string(config.ServiceName),
		Version: string(config.ServiceVersion),
		Deps:    deps,
	})
}

func (h *handler) GetSpecs(rw http.ResponseWriter, r *http.Request) {
	filename := "zookeeper_specs.yaml"
	rw.Header().Add("Content-Type", "application/openapi+yaml")
	rw.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	_, _ = rw.Write(openapi.ZookeeperSpec(h.selfURL))
}

func (h *handler) GetMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) SignIn1(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) AcrivateOrbisSocius(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) RequestOrbisSocius(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}
