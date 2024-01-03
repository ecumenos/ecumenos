package orbissocius

import "net/http"

type GetPingRespData struct {
	Ok bool `json:"ok"`
}

func (o *OrbisSocius) Health(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := o.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetPingRespData{Ok: true})
}

type GetInfoRespData struct {
	Name              string `json:"name"`
	PostgresIsRunning bool   `json:"postgres_is_running"`
}

func (o *OrbisSocius) Info(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writer := o.responseFactory.NewWriter(rw)
	writer.WriteSuccess(ctx, &GetInfoRespData{
		Name:              string(ServiceName),
		PostgresIsRunning: o.pg.Ping(ctx) == nil,
	})
}
