package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ecumenos/ecumenos/internal/docs"
	f "github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	gen "github.com/ecumenos/ecumenos/internal/generated/zookeeperadmin"
	"github.com/ecumenos/ecumenos/internal/openapi"
	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	responseFactory f.Factory
	service         *service.Service
	selfURL         string
	logger          *zap.Logger
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
		selfURL:         params.Config.AdminSelfURL,
		logger:          params.Logger,
	}
}

func (h *handler) auth(rw http.ResponseWriter, r *http.Request) context.Context {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)
	token, err := httputils.ExtractJWTBearerToken(r)
	if err != nil {
		_ = writer.WriteFail(ctx, nil, f.WithHTTPStatusCode(http.StatusUnauthorized), //nolint:errcheck
			f.WithCause(err), f.WithMessage("failed to get token"))
		h.logger.Error("can not extract JWT token from request", zap.Error(err))
		return nil
	}

	adminID, sessionID, err := h.service.AuthorizeAdmin(ctx, token)
	if err != nil {
		_ = writer.WriteFail(ctx, nil, f.WithHTTPStatusCode(http.StatusUnauthorized), //nolint:errcheck
			f.WithCause(err), f.WithMessage("failed to authorize"))
		h.logger.Error("can not authorize", zap.Error(err))
		return nil
	}
	ctx = contextutils.SetAdminID(ctx, adminID)
	ctx = contextutils.SetAdminSessionID(ctx, sessionID)

	return ctx
}

func (h *handler) GetDocs(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=UTF-8")
	_, _ = rw.Write(docs.ZookeeperAdminDocs(h.selfURL))
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
	filename := "zookeeper_admin_specs.yaml"
	rw.Header().Add("Content-Type", "application/openapi+yaml")
	rw.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	_, _ = rw.Write(openapi.ZookeeperAdminSpec(h.selfURL))
}

func (h *handler) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[gen.RefreshSessionRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	session, err := h.service.RefreshAdminSession(ctx, request.RefreshToken)
	if err != nil {
		_ = writer.WriteError(ctx, "failed refresh admin session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.RefreshSessionResponseData{ //nolint:errcheck
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		SessionId:    session.ID,
	})
}

func (h *handler) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[gen.SignInRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	if err := h.service.ValidateAdminCredentials(ctx, string(request.Email), request.Password); err != nil {
		_ = writer.WriteFail(ctx, "invalid email or password", f.WithCause(err), f.WithHTTPStatusCode(http.StatusUnauthorized)) //nolint:errcheck
		return
	}

	a, err := h.service.GetAdminByEmail(ctx, string(request.Email))
	if err != nil {
		_ = writer.WriteError(ctx, "can not get admin by email", err) //nolint:errcheck
		return
	}
	if a == nil {
		_ = writer.WriteFail(ctx, "invalid email", f.WithHTTPStatusCode(http.StatusNotFound)) //nolint:errcheck
		return
	}
	session, err := h.service.CreateAdminSession(ctx, a.ID)
	if err != nil {
		_ = writer.WriteError(ctx, "can not create admin session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.SignInResponseData{ //nolint:errcheck
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		SessionId:    session.ID,
	})
}

func (h *handler) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	sessionID, ok := contextutils.GetAdminSessionID(ctx)
	if !ok {
		_ = writer.WriteError(ctx, "something went wrong", errors.New("can not get session id from context")) //nolint:errcheck
		return
	}
	if err := h.service.DeleteAdminSession(ctx, sessionID); err != nil {
		_ = writer.WriteError(ctx, "failed detele admin session", err) //nolint:errcheck
		return
	}
	_ = writer.WriteSuccess(ctx, nil, f.WithHTTPStatusCode(http.StatusNoContent))
}
