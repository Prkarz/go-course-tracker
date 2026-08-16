package config

import "time"

// JWTSecret is the HMAC secret used to sign and verify JWT tokens.
// Consider overriding this via environment variables in production.
var JWTSecret = []byte("course-tracker-super-secret-key")

// JWTExpiryDuration is the default lifetime for generated JWTs.
var JWTExpiryDuration = 24 * time.Hour
