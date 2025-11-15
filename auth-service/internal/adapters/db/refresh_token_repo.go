package db

import (
	"context"
	"time"

	"github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/jwt"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, tokenHash string) (*jwt.Claims, error)
	Revoke(ctx context.Context, tokenHash string) error
	Rotate(ctx context.Context, oldHash, newHash string, newExpiresAt time.Time) error
}
