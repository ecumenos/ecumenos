package zookeeper

import (
	"net/http"
)

func (s *Server) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.zookeeper.Health())
}

func (s *Server) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := s.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, s.zookeeper.Info(ctx))
}
