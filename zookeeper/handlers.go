package zookeeper

import "net/http"

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (z *Zookeeper) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := z.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetPingRespData{Ok: true})
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (z *Zookeeper) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := z.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: z.pg.Ping(ctx) == nil,
	})
}
