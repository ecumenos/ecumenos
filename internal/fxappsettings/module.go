package fxappsettings

import (
	"context"
	"errors"

	"github.com/ecumenos/ecumenos/internal/fxappsettings/appsettings"
	"go.uber.org/fx"
)

type Config struct {
	RegionsPath string `json:"regionsPath"`
	LocalesPath string `json:"localesPath"`
}

var Module = fx.Options(
	fx.Provide(func(lc fx.Lifecycle, cfg *Config) (AppSettings, error) {
		if cfg.LocalesPath == "" || cfg.RegionsPath == "" {
			return nil, errors.New("locales path or regions path is empty")
		}

		m := appsettings.New()
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := m.LoadLocales(ctx, cfg.LocalesPath); err != nil {
					return err
				}

				return m.LoadRegions(ctx, cfg.RegionsPath)
			},
		})

		return m, nil
	}),
)

type AppSettings interface {
	GetCountries(onlyEnabled bool) []string
	GetRegionsByCountryCode(code string) []string
	ValidateRegionCode(code string) error
	ValidateCountryCode(code string) error
	GetLanguages(onlyEnabled bool) []string
	ValidateLanguageCode(code string) error
}
