package jwtmocktest

import (
	"bytes"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	jm "github.com/nayyara-cropsey/jwt-mock/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer()
	assert.NoError(t, err)

	defer server.Close()

	now := time.Now().Truncate(time.Second).UTC()
	exp := now.Add(time.Hour)
	token, err := server.GenerateJWT(jm.Claims{
		jwt.SubjectKey:    "olg387f",
		jwt.IssuedAtKey:   now.Unix(),
		jwt.ExpirationKey: exp.Unix(),
		jwt.IssuerKey:     "test",
		"email":           "vibrant_greider@xxx.com",
	})
	assert.NoError(t, err)

	jwsKeySet, err := jwk.Fetch(server.URL + "/.well-known/jwks.json")
	assert.NoError(t, err)

	parsedToken, err := jwt.Parse(bytes.NewReader([]byte(token)), jwt.WithKeySet(jwsKeySet))
	assert.NoError(t, err)

	assert.NoError(t, jwt.Verify(parsedToken, jwt.WithSubject("olg387f")))
	assert.NoError(t, jwt.Verify(parsedToken, jwt.WithIssuer("test")))
	assert.NoError(t, jwt.Verify(parsedToken, jwt.WithClaimValue(jwt.IssuedAtKey, now)))
	assert.NoError(t, jwt.Verify(parsedToken, jwt.WithClaimValue(jwt.ExpirationKey, exp)))
	assert.NoError(t, jwt.Verify(parsedToken, jwt.WithClaimValue("email", "vibrant_greider@xxx.com")))
	assert.NoError(t, jwt.Verify(parsedToken), jwt.WithKeySet(jwsKeySet))
}
