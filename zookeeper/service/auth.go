package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Authorization struct {
	JWTSigningKey []byte
}

func makeToken(subject string, scope string, exp time.Time) jwt.Token {
	tok := jwt.New()
	tok.Set("scope", scope)
	tok.Set("sub", subject)
	tok.Set("iat", time.Now().Unix())
	tok.Set("exp", exp.Unix())

	return tok
}

func (a *Authorization) CreateTokens(ctx context.Context, adminID int64, tokExp, refTokExp time.Time) (string, string, error) {
	accessTok := makeToken(fmt.Sprint(adminID), "access", tokExp)
	refreshTok := makeToken(fmt.Sprint(adminID), "refresh", refTokExp)

	rval := make([]byte, 10)
	rand.Read(rval)
	refreshTok.Set("jti", base64.StdEncoding.EncodeToString(rval))

	accSig, err := jwt.Sign(accessTok, jwt.WithKey(jwa.HS256, a.JWTSigningKey))
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}

	refSig, err := jwt.Sign(refreshTok, jwt.WithKey(jwa.HS256, a.JWTSigningKey))
	if err != nil {
		return "", "", fmt.Errorf("signing refresh token: %w", err)
	}

	return string(accSig), string(refSig), nil
}

func (a *Authorization) GetExpiredAt() (forToken, forRefreshToken time.Time) {
	return time.Now().Add(24 * time.Hour), time.Now().Add(7 * 24 * time.Hour)
}

func (a *Authorization) DecodeToken(token string) (jwt.Token, error) {
	t, err := jwt.ParseString(token, jwt.WithVerify(false), jwt.WithValidate(true))
	if err != nil {
		return nil, fmt.Errorf("decode token err: %w", err)
	}
	return t, nil
}
