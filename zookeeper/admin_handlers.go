package zookeeper

import (
	"net/http"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
)

func (s *AdminServer) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.zookeeper.Health())
}

func (s *AdminServer) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.zookeeper.Info(ctx))
}

func (s *AdminServer) SignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}

func (s *AdminServer) SignOut(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}

func (s *AdminServer) RefreshSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteError(ctx, "unimplemented", nil, fxresponsefactory.WithHTTPStatusCode(http.StatusNotImplemented))
}
