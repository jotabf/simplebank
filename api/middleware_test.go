package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jotabf/simplebank/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func Test_authMiddleware(t *testing.T) {
	type args struct {
		tokenMaker token.Maker
	}
	tests := []struct {
		name  string
		setup func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		check func(t *testing.T, rercoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			check: func(t *testing.T, rercoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rercoder.Code)
			},
		},
		{
			name:  "NoAuthorization",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			check: func(t *testing.T, rercoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rercoder.Code)
			},
		},
		{
			name: "UnsuportedAuthorization",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			check: func(t *testing.T, rercoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rercoder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
			},
			check: func(t *testing.T, rercoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rercoder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", 0*time.Nanosecond)
			},
			check: func(t *testing.T, rercoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rercoder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tt.setup(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tt.check(t, recorder)
		})
	}
}
