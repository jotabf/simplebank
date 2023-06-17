package token

import (
	"testing"
	"time"

	"github.com/jotabf/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestPasetoMakerKeySizeError(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(31))
	require.Error(t, err)
	require.EqualError(t, err, PasetoMakerInvalidSecretKeyError.Error())
	require.Nil(t, maker)
}

func TestPasetoMakerExpiredTokenError(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), time.Millisecond)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(2 * time.Millisecond)
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestPasetoMakerInvalidTokensError(t *testing.T) {
	maker1, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token1, err := maker1.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token1)

	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token1)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
