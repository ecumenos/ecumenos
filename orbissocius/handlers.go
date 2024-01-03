package orbissocius

import "net/http"

func (s *AdminServer) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.o.Health())
}

func (s *AdminServer) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.o.Info(ctx))
}

func (s *Server) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.o.Health())
}

func (s *Server) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.o.Info(ctx))
}
