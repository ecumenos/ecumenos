package zookeeper

import (
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
			ctx = contextutils.SetValue(ctx, contextutils.IPAddressKey, ip)
			ctx = contextutils.SetValue(ctx, contextutils.RequestIDKey, httputils.ExtractRequestID(r))
			ctx = contextutils.SetValue(ctx, contextutils.StartRequestTimestampKey, fmt.Sprint(time.Now().Unix()))

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

func NewAdminAuthorizationMiddleware(logger *zap.Logger, rf fxresponsefactory.Factory, z *Zookeeper) func(next http.Handler) http.Handler {
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
			adminID, session, err := z.Authorize(ctx, token)
			if err != nil {
				_ = writer.WriteFail(ctx, nil, fxresponsefactory.WithHTTPStatusCode(http.StatusUnauthorized),
					fxresponsefactory.WithCause(err), fxresponsefactory.WithMessage("failed to authorize")) //nolint:errcheck
				logger.Error("can not authorize", zap.Error(err))
				return
			}
			ctx = contextutils.SetValue(ctx, contextutils.AdminIDKey, fmt.Sprint(adminID))
			ctx = contextutils.SetValue(ctx, contextutils.AdminSessionIDKey, fmt.Sprint(session.ID))

			next.ServeHTTP(rw, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
