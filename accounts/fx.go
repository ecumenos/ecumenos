package accounts

import (
	"github.com/ecumenos/ecumenos/accounts/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Supply(config.ServiceName),
	fx.Supply(config.ServiceVersion),
)
