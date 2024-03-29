package httputils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ecumenos/ecumenos/internal/fxresponsefactory"
	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
	"github.com/ecumenos/ecumenos/internal/toolkit/netutils"
	"go.uber.org/zap"
)

func NewEnrichContextMiddleware(logger *zap.Logger, rf fxresponsefactory.Factory) func(next http.Handler) http.Handler {
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
			ctx = contextutils.SetIPAddress(ctx, ip)
			ctx = contextutils.SetRequestID(ctx, httputils.ExtractRequestID(r))
			ctx = contextutils.SetStartRequestTimestamp(ctx, time.Now())

			next.ServeHTTP(rw, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func NewRecoverMiddleware(logger *zap.Logger, rf fxresponsefactory.Factory) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			defer func() {
				if err := recover(); err != nil {
					_ = rf.NewWriter(rw).WriteError(ctx, "something went wrong", fmt.Errorf("unexpected error (err=%v)", err)) //nolint:errcheck
					logger.Error("recovering after panic", zap.Any("err", err))
					return
				}
			}()

			next.ServeHTTP(rw, r)
		}

		return http.HandlerFunc(fn)
	}
}

type Authorizer interface {
	Authorize(ctx context.Context, token string) (int64, int64, error)
}

func NewAdminAuthorizationMiddleware(logger *zap.Logger, rf fxresponsefactory.Factory, auth Authorizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			writer := rf.NewWriter(rw)

			token, err := httputils.ExtractJWTBearerToken(r)
			if err != nil {
				_ = writer.WriteFail(ctx, nil, fxresponsefactory.WithHTTPStatusCode(http.StatusUnauthorized),
					fxresponsefactory.WithCause(err), fxresponsefactory.WithMessage("failed to get token")) //nolint:errcheck
				logger.Error("can not extract JWT token from request", zap.Error(err))
				return
			}
			adminID, sessionID, err := auth.Authorize(ctx, token)
			if err != nil {
				_ = writer.WriteFail(ctx, nil, fxresponsefactory.WithHTTPStatusCode(http.StatusUnauthorized),
					fxresponsefactory.WithCause(err), fxresponsefactory.WithMessage("failed to authorize")) //nolint:errcheck
				logger.Error("can not authorize", zap.Error(err))
				return
			}
			ctx = contextutils.SetAdminID(ctx, adminID)
			ctx = contextutils.SetAdminSessionID(ctx, sessionID)

			next.ServeHTTP(rw, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
