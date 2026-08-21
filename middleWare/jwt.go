package middleWare

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Prkarz/course-tracker/config"
	"github.com/golang-jwt/jwt/v5"
)

func JWT_Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt_handler(w, r, next)
	})
}

func jwt_handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token_raw := r.Header.Get("Authorization")
	if token_raw == "" {
		writeError(w, http.StatusUnauthorized, "401_MISSING_TOKEN", "Authorization token is missing. Use the Bearer token format.")
		return
	}
	token_is_valid := strings.HasPrefix(token_raw, "Bearer ")
	if !token_is_valid {
		writeError(w, http.StatusUnauthorized, "401_INVALID_TOKEN_FORMAT", "Invalid token format. Use Bearer <token>.")
		return
	}
	token_trim := strings.TrimPrefix(token_raw, "Bearer ")

	token, err := jwt.Parse(token_trim, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.GetJWTSecret(), nil
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, "401_TOKEN_PARSE_FAILED", "Failed to validate the token. It may be expired or invalid.")
		return
	}
	if !token.Valid {
		writeError(w, http.StatusUnauthorized, "401_TOKEN_INVALID", "The token is invalid or expired. Please log in again.")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		writeError(w, http.StatusUnauthorized, "401_CLAIMS_INVALID", "Token claims are invalid or malformed.")
		return
	}
	floatID, ok := claims["user_id"].(float64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "401_CLAIMS_INVALID", "The token user ID is missing or invalid.")
		return
	}
	userID := int(floatID)
	ctx := context.WithValue(r.Context(), "userID", userID)
	reqWithContext := r.WithContext(ctx)
	next(w, reqWithContext)
}
