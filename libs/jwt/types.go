package jwt

import "time"

// minimal DTOs and config for the JWT service

// UserInfo is the minimal user information used to generate tokens.
// Keep it small and stable — convert domain entities to this type before calling JWT service.
type UserInfo struct {
    UserID uint64
    Email  string
    Role   string
}

// Claims is the subset we return to application after validating an access token.
type Claims struct {
    UserID uint64 `json:"uid"`
    Email  string `json:"email,omitempty"`
    Role   string `json:"role,omitempty"`
}

// Config for the JWT/RSA service
type Config struct {
    Issuer     string        // iss claim
    Audience   string        // aud claim
    AccessTTL  time.Duration // how long access JWT is valid
    RefreshTTL time.Duration // how long refresh token stays valid
    KeyID      string        // kid for the signing key
}
