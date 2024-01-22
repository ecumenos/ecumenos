package zookeeper

import (
	"github.com/ecumenos/ecumenos/zookeeper/admin"
	"github.com/ecumenos/ecumenos/zookeeper/app"
	"github.com/ecumenos/ecumenos/zookeeper/config"
	"github.com/ecumenos/ecumenos/zookeeper/repository"
	"github.com/ecumenos/ecumenos/zookeeper/service"
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
