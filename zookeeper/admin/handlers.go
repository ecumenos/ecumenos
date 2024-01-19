package admin

import (
	"net/http"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
)

type GetHealthRespData struct {
	Ok bool `json:"ok"`
}

func (s *Server) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, GetHealthRespData{Ok: true})
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (s *Server) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := s.service.PingServices(ctx)
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, GetInfoRespData{
		Name:              string(s.serviceName),
		PostgresIsRunning: result.PostgresIsRunning,
	})
}

func (s *Server) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}

func (s *Server) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}

func (s *Server) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}
