package httputils

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

func RobustHTTPClient(logger *zap.Logger) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.Logger = retryablehttp.LeveledLogger(LeveledZap{inner: logger})
	retryClient.CheckRetry = RetryPolicy
	client := retryClient.StandardClient()
	client.Timeout = 30 * time.Second
	return client
}

type LeveledZap struct {
	inner *zap.Logger
}

// re-writes HTTP client ERROR to WARN level (because of retries)
func (l LeveledZap) Error(msg string, keysAndValues ...interface{}) {
	l.inner.Sugar().Warnw(msg, keysAndValues...)
}

func (l LeveledZap) Warn(msg string, keysAndValues ...interface{}) {
	l.inner.Sugar().Warnw(msg, keysAndValues...)
}

func (l LeveledZap) Info(msg string, keysAndValues ...interface{}) {
	l.inner.Sugar().Infow(msg, keysAndValues...)
}

func (l LeveledZap) Debug(msg string, keysAndValues ...interface{}) {
	l.inner.Sugar().Debugw(msg, keysAndValues...)
}

// RetryPolicy is a custom wrapper around retryablehttp.DefaultRetryPolicy.
// It treats `429 Too Many Requests` as non-retryable, so the application can decide
// how to deal with rate-limiting.
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if err == nil && resp.StatusCode == http.StatusTooManyRequests {
		return false, nil
	}
	// TODO: implement returning errors on non-200 responses w/o introducing circular dependencies.
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
