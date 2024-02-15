package service_test

import (
	"context"
	"testing"

	"github.com/ecumenos/ecumenos/internal/toolkit/primitives"
	"github.com/ecumenos/ecumenos/zookeeper/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	a := service.Authorization{JWTSigningKey: []byte("qwerty")}
	ctx := context.Background()
	adminID := int64(1234567890)
	expTok, expRefTok := a.GetExpiredAt()
	tok, refTok, err := a.CreateTokens(ctx, adminID, expTok, expRefTok)
	require.NoError(t, err)

	jwtTok, err := a.DecodeToken(tok)
	require.NoError(t, err)
	subject := jwtTok.Subject()
	actualAdminID, err := primitives.StringToInt64(subject)
	require.NoError(t, err)
	assert.Equal(t, adminID, actualAdminID)

	jwtRefTok, err := a.DecodeToken(refTok)
	require.NoError(t, err)
	refSubject := jwtRefTok.Subject()
	actualRefAdminID, err := primitives.StringToInt64(refSubject)
	require.NoError(t, err)
	assert.Equal(t, adminID, actualRefAdminID)
}
