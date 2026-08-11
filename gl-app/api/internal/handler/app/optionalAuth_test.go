package app

import (
	"net/http/httptest"
	"testing"
	"time"

	"gl-app/pkg/ctxdata"

	"github.com/golang-jwt/jwt/v4"
)

func TestOptionalJWTUserID(t *testing.T) {
	const secret = "test-secret"
	validToken := signedTestToken(t, secret, time.Now().Add(time.Hour).Unix(), 42)
	expiredToken := signedTestToken(t, secret, time.Now().Add(-time.Hour).Unix(), 42)

	tests := []struct {
		name   string
		header string
		want   int64
	}{
		{name: "guest", want: 0},
		{name: "valid token", header: "Bearer " + validToken, want: 42},
		{name: "expired token", header: "Bearer " + expiredToken, want: 0},
		{name: "invalid token", header: "Bearer invalid", want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/app/v1/users/42/works", nil)
			if test.header != "" {
				request.Header.Set("Authorization", test.header)
			}
			if got := optionalJWTUserID(request, secret); got != test.want {
				t.Fatalf("optionalJWTUserID() = %d, want %d", got, test.want)
			}
		})
	}
}

func signedTestToken(t *testing.T, secret string, expiresAt, userID int64) string {
	t.Helper()
	claims := jwt.MapClaims{
		"exp":                   expiresAt,
		ctxdata.CtxKeyJwtUserId: userID,
	}
	value, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return value
}
