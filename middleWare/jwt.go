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
		http.Error(w, "Empty token", http.StatusUnauthorized)
		return
	}
	token_is_valid := strings.HasPrefix(token_raw, "Bearer ")
	if !token_is_valid {
		http.Error(w, "Token not proper", http.StatusUnauthorized)
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
		http.Error(w, "Unable to parse token", http.StatusUnauthorized)
		return
	}
	if !token.Valid {
		http.Error(w, "Token in invalid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "CLaims in invalid", http.StatusUnauthorized)
		return
	}
	floatID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Claims in invalid", http.StatusUnauthorized)
		return
	}
	userID := int(floatID)
	ctx := context.WithValue(r.Context(), "userID", userID)
	reqWithContext := r.WithContext(ctx)
	next(w, reqWithContext)
}
