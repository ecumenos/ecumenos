package fxresponsefactory

import (
	"net/http"

	"github.com/ecumenos/ecumenos/internal/fxtypes"
	"go.uber.org/zap"
)

type Factory interface {
	NewWriter(rw http.ResponseWriter) Writer
}

type Config struct {
	WriteLogs bool
}

func NewFactory(l *zap.Logger, cfg *Config, version fxtypes.Version) Factory {
	return &factory{
		l:         l,
		writeLogs: cfg.WriteLogs,
		version:   version,
	}
}

type factory struct {
	l         *zap.Logger
	writeLogs bool
	version   fxtypes.Version
}

func (f *factory) NewWriter(rw http.ResponseWriter) Writer {
	return NewWriter(f.l, rw, f.version, f.writeLogs)
}
