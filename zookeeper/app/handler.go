package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ecumenos/ecumenos/internal/docs"
	f "github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	gen "github.com/ecumenos/ecumenos/internal/generated/zookeeper"
	"github.com/ecumenos/ecumenos/internal/openapi"
	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
	models "github.com/ecumenos/ecumenos/models/zookeeper"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	"github.com/oapi-codegen/runtime/types"
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
		selfURL:         params.Config.AppSelfURL,
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

	comptusID, sessionID, err := h.service.AuthorizeComptus(ctx, token)
	if err != nil {
		_ = writer.WriteFail(ctx, nil, f.WithHTTPStatusCode(http.StatusUnauthorized), //nolint:errcheck
			f.WithCause(err), f.WithMessage("failed to authorize"))
		h.logger.Error("can not authorize", zap.Error(err))
		return nil
	}
	ctx = contextutils.SetComptusID(ctx, comptusID)
	ctx = contextutils.SetComptusSessionID(ctx, sessionID)

	return ctx
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

func mapModelComptusToGenComptus(v *models.Comptus) gen.Comptus {
	return gen.Comptus{
		Id:       v.ID,
		Country:  v.Patria,
		Email:    types.Email(v.Email),
		Language: v.Lingua,
	}
}

func (h *handler) GetMe(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	comptusID, ok := contextutils.GetComptusID(ctx)
	if !ok {
		_ = writer.WriteError(ctx, "something went wrong", errors.New("can not get comptus id from context")) //nolint:errcheck
		return
	}
	c, err := h.service.GetComptusByID(ctx, comptusID)
	if err != nil {
		_ = writer.WriteError(ctx, "failed get comptus", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, mapModelComptusToGenComptus(c)) //nolint:errcheck
}

func (h *handler) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[gen.RefreshSessionRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	session, err := h.service.RefreshComptusSession(ctx, request.RefreshToken)
	if err != nil {
		_ = writer.WriteError(ctx, "failed refresh comptus session", err) //nolint:errcheck
		return
	}

	c, err := h.service.GetComptusByID(ctx, session.ComptusID)
	if err != nil {
		_ = writer.WriteError(ctx, "failed refresh comptus session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.RefreshSessionResponseData{ //nolint:errcheck
		Self: mapModelComptusToGenComptus(c),
		Tokens: gen.AuthTokenPair{
			RefreshToken: session.RefreshToken,
			Token:        session.Token,
		},
		SessionId: session.ID,
	})
}

func (h *handler) SignUp(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[gen.SignUpRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}

	c, err := h.service.CreateComptus(ctx, string(request.Email), request.Password, string(request.Country), string(request.Language))
	if err != nil {
		_ = writer.WriteError(ctx, "can not create comptus", err) //nolint:errcheck
		return
	}
	session, err := h.service.CreateComptusSession(ctx, c.ID)
	if err != nil {
		_ = writer.WriteError(ctx, "can not create comptus session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.SignUpResponseData{ //nolint:errcheck
		Self: mapModelComptusToGenComptus(c),
		Tokens: gen.AuthTokenPair{
			RefreshToken: session.RefreshToken,
			Token:        session.Token,
		},
		SessionId: session.ID,
	})
}

func (h *handler) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	sessionID, ok := contextutils.GetComptusSessionID(ctx)
	if !ok {
		_ = writer.WriteError(ctx, "something went wrong", errors.New("can not get session id from context")) //nolint:errcheck
		return
	}
	if err := h.service.DeleteComptusSession(ctx, sessionID); err != nil {
		_ = writer.WriteError(ctx, "failed detele comptus session", err) //nolint:errcheck
		return
	}
	_ = writer.WriteSuccess(ctx, nil, f.WithHTTPStatusCode(http.StatusNoContent))
}

func (h *handler) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := h.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[gen.SignInRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	if err := h.service.ValidateComptusCredentials(ctx, string(request.Email), request.Password); err != nil {
		_ = writer.WriteFail(ctx, "invalid email or password", f.WithCause(err), f.WithHTTPStatusCode(http.StatusUnauthorized)) //nolint:errcheck
		return
	}

	c, err := h.service.GetComptusByEmail(ctx, string(request.Email))
	if err != nil {
		_ = writer.WriteError(ctx, "can not get comptus by email", err) //nolint:errcheck
		return
	}
	if c == nil {
		_ = writer.WriteFail(ctx, "invalid email", f.WithHTTPStatusCode(http.StatusNotFound)) //nolint:errcheck
		return
	}
	session, err := h.service.CreateComptusSession(ctx, c.ID)
	if err != nil {
		_ = writer.WriteError(ctx, "can not create comptus session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.SignInResponseData{ //nolint:errcheck
		Self: mapModelComptusToGenComptus(c),
		Tokens: gen.AuthTokenPair{
			RefreshToken: session.RefreshToken,
			Token:        session.Token,
		},
		SessionId: session.ID,
	})
}

func (h *handler) GetCountries(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	countries := h.service.GetOrbisSociusCountries()

	_ = writer.WriteSuccess(ctx, gen.CountriesResponseData(countries)) //nolint:errcheck
}

func (h *handler) GetCountryRegions(rw http.ResponseWriter, r *http.Request, countryCode string) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	regions := h.service.GetOrbisSociusRegions(countryCode)

	_ = writer.WriteSuccess(ctx, gen.CountryRegionsResponseData(regions)) //nolint:errcheck
}

func (h *handler) GetLanguages(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	languages := h.service.GetOrbisSociusLanguages()

	_ = writer.WriteSuccess(ctx, gen.LanguagesResponseData(languages)) //nolint:errcheck
}

func (h *handler) ActivateOrbisSocius(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}

	writer := h.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", nil, f.WithHTTPStatusCode(http.StatusNotImplemented)) //nolint:errcheck
}

func (h *handler) RequestOrbisSocius(rw http.ResponseWriter, r *http.Request) {
	ctx := h.auth(rw, r)
	if ctx == nil {
		return
	}
	writer := h.responseFactory.NewWriter(rw)
	comptusID, ok := contextutils.GetComptusID(ctx)
	if !ok {
		_ = writer.WriteError(ctx, "something went wrong", errors.New("can not get comptus id from context")) //nolint:errcheck
		return
	}

	request, err := httputils.DecodeBody[gen.RequestOrbisSociusRequest](h.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}

	if _, err := h.service.MakeCreateOrbisSociusLaunchRequest(ctx, comptusID, request.Region, request.Name, request.Description, request.Url); err != nil {
		_ = writer.WriteFail(ctx, "can not create orbis socius launch request", f.WithCause(err)) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, gen.RequestOrbisSociusResponseData{ //nolint:errcheck
		Ok: true,
	})
}
