package admin

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHandler),
	fx.Provide(NewServer),
)
