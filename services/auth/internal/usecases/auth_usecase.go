package usecases

import (
	"context"
	"errors"

	"bitka/auth-service/internal/adapters/db"
	"bitka/auth-service/internal/entities"
	errs "bitka/common"
	"bitka/jwt"
	"bitka/auth-service/pkg/passhash"
	"github.com/rs/zerolog"
)

type AuthUsecase interface {
	Register(ctx context.Context, email, password string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
}

type authUsecase struct {
	userRepo db.UserRepository
	jwtSvc   jwt.Service
	hasher   passhash.Service
	logger   *zerolog.Logger
}

func NewAuthUsecase(u db.UserRepository, j jwt.Service, h passhash.Service, logger *zerolog.Logger) AuthUsecase {
	return &authUsecase{userRepo: u, jwtSvc: j, hasher: h, logger: logger}
}

func (a *authUsecase) Register(ctx context.Context, email, password string) (*entities.User, error) {
	a.logger.Debug().Str("email", email).Msg("Register start")

	existing, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// If repo returned CodedError we just bubble it (maybe wrap)
		a.logger.Error().Err(err).Str("email", email).Msg("GetByEmail failed")
		return nil, errs.Wrap(err, errs.CodeInternal, "failed to check existing user")
	}
	if existing != nil {
		a.logger.Info().Str("email", email).Msg("register conflict")
		return nil, errs.New(errs.CodeConflict, "email already registered")
	}

	hashed, err := a.hasher.Hash(password)
	if err != nil {
		a.logger.Error().Err(err).Msg("hash failed")
		return nil, errs.Wrap(err, errs.CodeInternal, "hash failed")
	}

	user := &entities.User{Email: email, Password: hashed /* ... */}
	created, err := a.userRepo.Create(ctx, user)
	if err != nil {
		a.logger.Error().Err(err).Str("email", email).Msg("create user failed")
		return nil, err // already wrapped in repo
	}

	a.logger.Info().Uint64("user_id", created.ID).Str("email", created.Email).Msg("register success")
	return created, nil
}

// TODO: add logging similar to Register
func (a *authUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if u == nil {
		return "", "", errors.New("invalid credentials")
	}
	if !a.hasher.Verify(u.Password, password) {
		return "", "", errors.New("invalid credentials")
	}
	// create tokens
	access, err := a.jwtSvc.GenerateAccessToken(ctx, u)
	if err != nil {
		return "", "", err
	}
	refresh, _, _, err := a.jwtSvc.GenerateRefreshToken(ctx, u)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
