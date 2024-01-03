package pds

import "net/http"

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (pds *PDS) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := pds.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetPingRespData{Ok: true})
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (pds *PDS) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := pds.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: pds.pg.Ping(ctx) == nil,
	})
}
