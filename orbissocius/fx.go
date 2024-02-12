package orbissocius

import (
	"github.com/ecumenos/ecumenos/orbissocius/admin"
	"github.com/ecumenos/ecumenos/orbissocius/app"
	"github.com/ecumenos/ecumenos/orbissocius/config"
	"github.com/ecumenos/ecumenos/orbissocius/repository"
	"github.com/ecumenos/ecumenos/orbissocius/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	repository.Module,
	service.Module,
	app.Module,
	admin.Module,
	fx.Supply(config.ServiceName),
	fx.Supply(config.ServiceVersion),
	fx.Provide(
		NewMigrationsRunner,
	),
)
