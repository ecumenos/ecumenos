package httputils

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ecumenos/ecumenos/internal/toolkit/contextutils"
	"github.com/ecumenos/ecumenos/internal/toolkit/primitives"
	"github.com/ecumenos/ecumenos/internal/toolkit/random"
)

func ExtractRequestID(r *http.Request) string {
	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		return reqID
	}

	return random.GenUUIDString()
}

func GetRequestDuration(ctx context.Context) (time.Duration, error) {
	str := contextutils.GetValueFromContext(ctx, contextutils.StartRequestTimestampKey)
	if str == "" {
		return 0, nil
	}
	start, err := primitives.StringToInt64(str)
	if err != nil {
		return 0, err
	}
	diff := time.Now().Unix() - start

	return time.Duration(diff) * time.Second, nil
}

func ExtractJWTBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("token prefix is missing")
	}

	return authHeaderParts[1], nil
}
