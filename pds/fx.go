package pds

import (
	"github.com/ecumenos/ecumenos/pds/admin"
	"github.com/ecumenos/ecumenos/pds/app"
	"github.com/ecumenos/ecumenos/pds/config"
	"github.com/ecumenos/ecumenos/pds/repository"
	"github.com/ecumenos/ecumenos/pds/service"
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
