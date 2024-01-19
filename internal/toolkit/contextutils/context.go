package contextutils

import "context"

type ContextKey string

const (
	RequestIDKey             ContextKey = "k_request_id"
	StartRequestTimestampKey ContextKey = "k_start_request_timestamp"
	IPAddressKey             ContextKey = "k_ip_address"
	AdminIDKey               ContextKey = "k_admin_id"
	AccountIDKey             ContextKey = "k_account_id"
	AdminSessionIDKey        ContextKey = "k_admin_session_id"
	SessionIDKey             ContextKey = "k_session_id"
)

func GetValueFromContext(ctx context.Context, key ContextKey) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}

	return ""
}

func SetValue(ctx context.Context, key ContextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
