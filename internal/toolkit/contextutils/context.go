package contextutils

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type ctxKey string

const (
	requestIDKey             ctxKey = "k_request_id"
	startRequestTimestampKey ctxKey = "k_start_request_timestamp"
	ipAddressKey             ctxKey = "k_ip_address"
	adminIDKey               ctxKey = "k_admin_id"
	accountIDKey             ctxKey = "k_account_id"
	adminSessionIDKey        ctxKey = "k_admin_session_id"
	sessionIDKey             ctxKey = "k_session_id"
)

func getValueFromContext(ctx context.Context, key ctxKey) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}

	return ""
}

func setValue(ctx context.Context, key ctxKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

func SetRequestID(ctx context.Context, v string) context.Context {
	return setValue(ctx, requestIDKey, v)
}

func GetRequestID(ctx context.Context) string {
	return getValueFromContext(ctx, requestIDKey)
}

func SetStartRequestTimestamp(ctx context.Context, v time.Time) context.Context {
	return setValue(ctx, startRequestTimestampKey, fmt.Sprint(v.UnixMilli()))
}

func GetStartRequestTimestamp(ctx context.Context) (time.Time, bool) {
	v := getValueFromContext(ctx, startRequestTimestampKey)
	if v == "" {
		return time.Time{}, false
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	if i == 0 {
		return time.Time{}, false
	}

	return time.UnixMilli(i), true
}

func SetIPAddress(ctx context.Context, v string) context.Context {
	return setValue(ctx, ipAddressKey, v)
}

func GetIPAddress(ctx context.Context) string {
	return getValueFromContext(ctx, ipAddressKey)
}

func SetAdminID(ctx context.Context, v int64) context.Context {
	return setValue(ctx, adminIDKey, fmt.Sprint(v))
}

func GetAdminID(ctx context.Context) (int64, bool) {
	v := getValueFromContext(ctx, adminIDKey)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func SetAccountID(ctx context.Context, v int64) context.Context {
	return setValue(ctx, accountIDKey, fmt.Sprint(v))
}

func GetAccountID(ctx context.Context) (int64, bool) {
	v := getValueFromContext(ctx, accountIDKey)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func SetAdminSessionID(ctx context.Context, v int64) context.Context {
	return setValue(ctx, adminSessionIDKey, fmt.Sprint(v))
}

func GetAdminSessionID(ctx context.Context) (int64, bool) {
	v := getValueFromContext(ctx, adminSessionIDKey)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func SetSessionID(ctx context.Context, v int64) context.Context {
	return setValue(ctx, sessionIDKey, fmt.Sprint(v))
}

func GetSessionID(ctx context.Context) (int64, bool) {
	v := getValueFromContext(ctx, sessionIDKey)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}
