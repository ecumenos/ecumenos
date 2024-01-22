package admin

import (
	"errors"
	"net/http"

	f "github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
)

type GetHealthRespData struct {
	Ok bool `json:"ok"`
}

func (s *Server) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, GetHealthRespData{Ok: true}) //nolint:errcheck
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (s *Server) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := s.service.PingServices(ctx)
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, GetInfoRespData{ //nolint:errcheck
		Name:              string(s.serviceName),
		PostgresIsRunning: result.PostgresIsRunning,
	})
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRespData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    int64  `json:"session_id"`
}

func (s *Server) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[SignInReq](s.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	if err := s.service.ValidateAdminCredentials(ctx, request.Email, request.Password); err != nil {
		_ = writer.WriteFail(ctx, "invalid email or password", f.WithCause(err), f.WithHTTPStatusCode(http.StatusUnauthorized)) //nolint:errcheck
		return
	}

	a, err := s.service.GetAdminByEmail(ctx, request.Email)
	if err != nil {
		_ = writer.WriteError(ctx, "can not get admin by email", err) //nolint:errcheck
		return
	}
	if a == nil {
		_ = writer.WriteFail(ctx, "invalid email", f.WithHTTPStatusCode(http.StatusNotFound)) //nolint:errcheck
		return
	}
	session, err := s.service.CreateAdminSession(ctx, a.ID)
	if err != nil {
		_ = writer.WriteError(ctx, "can not create admin session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, SignInRespData{ //nolint:errcheck
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		SessionID:    session.ID,
	})
}

type SignOutRespData struct {
	Ok bool `json:"ok"`
}

func (s *Server) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	sessionID, ok := contextutils.GetAdminSessionID(ctx)
	if !ok {
		_ = writer.WriteError(ctx, "something went wrong", errors.New("can not get session id from context")) //nolint:errcheck
		return
	}
	if err := s.service.DeleteAdminSession(ctx, sessionID); err != nil {
		_ = writer.WriteError(ctx, "failed detele admin session", err) //nolint:errcheck
		return
	}
	_ = writer.WriteSuccess(ctx, SignOutRespData{ //nolint:errcheck
		Ok: true,
	})
}

type RefreshSessionReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshSessionRespData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    int64  `json:"session_id"`
}

func (s *Server) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)

	request, err := httputils.DecodeBody[RefreshSessionReq](s.logger, r)
	if err != nil {
		_ = writer.WriteFail(ctx, "invalid body", f.WithCause(err)) //nolint:errcheck
		return
	}
	session, err := s.service.RefreshAdminSession(ctx, request.RefreshToken)
	if err != nil {
		_ = writer.WriteError(ctx, "failed refresh admin session", err) //nolint:errcheck
		return
	}

	_ = writer.WriteSuccess(ctx, RefreshSessionRespData{ //nolint:errcheck
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		SessionID:    session.ID,
	})
}
