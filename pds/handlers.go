package pds

import "net/http"

func (s *AdminServer) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, s.pds.Health()) //nolint:errcheck
}

func (s *AdminServer) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, s.pds.Info(ctx)) //nolint:errcheck
}

func (s *Server) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, s.pds.Health()) //nolint:errcheck
}

func (s *Server) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	_ = writer.WriteSuccess(ctx, s.pds.Info(ctx)) //nolint:errcheck
}
