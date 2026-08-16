package middleWare

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Prkarz/course-tracker/config"
)

func JWT_Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt_handler(w, r, next)
	})
}

func jwt_handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token_raw := r.Header.Get("Authorization")
	if token_raw == "" {
		http.Error(w, "[401_MISSING_TOKEN] Authorization token is missing. Include 'Authorization: Bearer <token>' in request headers.", http.StatusUnauthorized)
		return
	}
	token_is_valid := strings.HasPrefix(token_raw, "Bearer ")
	if !token_is_valid {
		http.Error(w, "[401_INVALID_TOKEN_FORMAT] Invalid token format. Must use 'Bearer <token>' format.", http.StatusUnauthorized)
		return
	}
	token_trim := strings.TrimPrefix(token_raw, "Bearer ")

	token, err := jwt.Parse(token_trim, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing Method")
		}
		return config.JWTSecret, nil
	})
	if err != nil {
		http.Error(w, "[401_TOKEN_PARSE_FAILED] Failed to parse or validate token. Token may be expired or tampered.", http.StatusUnauthorized)
		return
	}
	if !token.Valid {
		http.Error(w, "[401_TOKEN_INVALID] The provided token is invalid or has expired. Please log in again.", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "[401_CLAIMS_INVALID] Token claims are invalid or malformed.", http.StatusUnauthorized)
		return
	}
	floatID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "[401_CLAIMS_INVALID] User ID claim is missing or invalid from token.", http.StatusUnauthorized)
		return
	}
	userID := int(floatID)
	ctx := context.WithValue(r.Context(), "userID", userID)
	reqWithContext := r.WithContext(ctx)
	next(w, reqWithContext)
}
