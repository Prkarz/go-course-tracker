package config

import (
	"os"
	"time"
)

// JWTSecret is the HMAC secret used to sign and verify JWT tokens.
// Consider overriding this via environment variables in production.

// JWTExpiryDuration is the default lifetime for generated JWTs.
var JWTExpiryDuration = 24 * time.Hour

func GetJWTSecret() []byte {
	JWTSecret := []byte(os.Getenv("SECRET_SAUCE"))

	if string(JWTSecret) == "" {
		// Fallback or panic if the secret is missing in production
		return []byte("fallback_development_secret")
	}
	return JWTSecret
}
