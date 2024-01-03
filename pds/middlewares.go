package pds

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/go-toolkit/contextutils"
	"github.com/ecumenos/go-toolkit/httputils"
	"github.com/ecumenos/go-toolkit/netutils"
	"go.uber.org/zap"
)

func NewEnrichContextMiddleware(logger *zap.Logger, rf fxrf.Factory) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			writer := rf.NewWriter(rw)

			ip, err := netutils.ExtractIPAddress(r)
			if err != nil {
				_ = writer.WriteFail(ctx, err, nil) //nolint:errcheck
				logger.Error("can not extract IP address from request", zap.Error(err))
				return
			}
			ctx = contextutils.SetValue(ctx, contextutils.IPAddressKey, ip)
			ctx = contextutils.SetValue(ctx, contextutils.RequestIDKey, httputils.ExtractRequestID(r))
			ctx = contextutils.SetValue(ctx, contextutils.StartRequestTimestampKey, fmt.Sprint(time.Now().UnixNano()))

			next.ServeHTTP(rw, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func NewRecoverMiddleware(logger *zap.Logger, rf fxrf.Factory) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			defer func() {
				if err := recover(); err != nil {
					_ = rf.NewWriter(rw).WriteError(ctx, "something went wrong", fmt.Errorf("unexpected error (err=%v)", err)) //nolint:errcheck
					logger.Error("can not get request duration", zap.Any("err", err))
					return
				}
			}()

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}
