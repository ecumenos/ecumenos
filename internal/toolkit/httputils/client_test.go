package httputils_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ecumenos/ecumenos/internal/toolkit/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRobustHTTPClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const url = "https://httpbin.org/ip"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)

	l := zaptest.NewLogger(t)
	c := httputils.RobustHTTPClient(l)
	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
