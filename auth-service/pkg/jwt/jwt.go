package jwt

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/entities"
)

// Service is the interface used by usecases / middleware
type Service interface {
	GenerateAccessToken(ctx context.Context, u *entities.User) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)

	// Refresh token flow
	GenerateRefreshToken(ctx context.Context, u *entities.User) (plain string, hashed string, expiresAt time.Time, err error)
	// ValidateRefreshToken takes hashed token (from DB) + provided plain token and validates expiry & equality
	// We provide helper to HashRefreshToken and CompareRefreshToken
	HashRefreshToken(plain string) (string, error)
	// Utility: random opaque token (plain)
	NewOpaqueToken(nBytes int) (string, error)
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type HMACConfig struct {
	AccessSecret []byte
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	Issuer       string
	Audience     string
}

type hmacService struct {
	cfg HMACConfig
}

// NewHMACService returns a Service using HS256 signing
func NewHMACService(cfg HMACConfig) Service {
	return &hmacService{cfg: cfg}
}

func (s *hmacService) GenerateAccessToken(ctx context.Context, u *entities.User) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    s.cfg.Issuer,
		Subject:   "", // optional; we'll use custom fields
		Audience:  jwt.ClaimStrings{s.cfg.Audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        "", // optional jti
	}

	custom := map[string]any{
		"uid":   u.ID,
		"email": u.Email,
		"role":  u.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   claims.Issuer,
		"aud":   claims.Audience,
		"exp":   claims.ExpiresAt.Unix(),
		"iat":   claims.IssuedAt.Unix(),
		"uid":   custom["uid"],
		"email": custom["email"],
		"role":  custom["role"],
	})

	signed, err := token.SignedString(s.cfg.AccessSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *hmacService) ValidateAccessToken(ctx context.Context, tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.cfg.AccessSecret, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		// check for expired specifically
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// extract claims safely
	uidFloat, _ := claimsMap["uid"].(float64)
	uid := uint64(uidFloat)
	email, _ := claimsMap["email"].(string)
	role, _ := claimsMap["role"].(string)

	return &Claims{UserID: uid, Email: email, Role: role}, nil
}

// ---- Refresh token helpers ----
// NewOpaqueToken: create a URL-safe base64 opaque token
func (s *hmacService) NewOpaqueToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashRefreshToken: hash token for DB storage (use SHA256)
func (s *hmacService) HashRefreshToken(plain string) (string, error) {
	h := sha256.Sum256([]byte(plain))
	return base64.RawURLEncoding.EncodeToString(h[:]), nil
}

// GenerateRefreshToken: create plain token and hashed token for DB storing
func (s *hmacService) GenerateRefreshToken(ctx context.Context, u *entities.User) (plain string, hashed string, expiresAt time.Time, err error) {
	plain, err = s.NewOpaqueToken(32) // 32 bytes -> 43 chars base64; adjust
	if err != nil {
		return "", "", time.Time{}, err
	}
	hashed, err = s.HashRefreshToken(plain)
	if err != nil {
		return "", "", time.Time{}, err
	}
	expiresAt = time.Now().Add(s.cfg.RefreshTTL).UTC()
	return plain, hashed, expiresAt, nil
}
