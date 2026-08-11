package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gl-app/pkg/ctxdata"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/token"
)

var optionalTokenParser = token.NewTokenParser()

// optionalJWTUserID returns a verified user ID when a public request carries a
// valid Bearer token. Missing, expired, or invalid tokens are treated as a
// guest request, so public endpoints remain accessible without authentication.
func optionalJWTUserID(r *http.Request, secret string) int64 {
	if r == nil || strings.TrimSpace(r.Header.Get("Authorization")) == "" || secret == "" {
		return 0
	}
	parsed, err := optionalTokenParser.ParseToken(r, secret, "")
	if err != nil || parsed == nil || !parsed.Valid {
		return 0
	}
	if _, ok := parsed.Method.(*jwt.SigningMethodHMAC); !ok {
		return 0
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}
	return positiveInt64Claim(claims[ctxdata.CtxKeyJwtUserId])
}

func positiveInt64Claim(value any) int64 {
	var userID int64
	switch typed := value.(type) {
	case json.Number:
		userID, _ = typed.Int64()
	case float64:
		userID = int64(typed)
	case string:
		userID, _ = strconv.ParseInt(typed, 10, 64)
	}
	if userID <= 0 {
		return 0
	}
	return userID
}
