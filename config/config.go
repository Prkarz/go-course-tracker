package config

// JWTSecret is the HMAC secret used to sign and verify JWT tokens.
// Consider overriding this via environment variables in production.
var JWTSecret = []byte("course-tracker-super-secret-key")
