package app

import (
	"errors"
	"net/http"

	f "github.com/ecumenos/ecumenos/internal/fxresponsefactory"
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

type RequestOrbisSociusLaunchReq struct {
	Email string `json:"email"`
}

func (s *Server) RequestOrbisSociusLaunch(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteError(ctx, "not implemented", errors.New("not implemented"), f.WithHTTPStatusCode(http.StatusNotImplemented))
}
