package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/entities"
	errs "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type GormUser struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null;default:'user'"`
	IsActive  bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GormUser) TableName() string { return "auth.users" } // e.g. schema auth

// conversion helpers
func (g *GormUser) ToEntity() *entities.User {
	return &entities.User{
		ID:        g.ID,
		Email:     g.Email,
		Password:  g.Password,
		Role:      g.Role,
		IsActive:  g.IsActive,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

func FromEntity(u *entities.User) *GormUser {
	return &GormUser{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserRepository interface {
	Create(ctx context.Context, u *entities.User) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id uint64) (*entities.User, error)
	Update(ctx context.Context, u *entities.User) error
}

type userRepoImpl struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewUserRepo(db *gorm.DB, logger *zerolog.Logger) UserRepository {
	return &userRepoImpl{
		db:     db,
		logger: logger,
	}
}

func (r *userRepoImpl) Create(ctx context.Context, u *entities.User) (*entities.User, error) {
	gu := FromEntity(u)
	if err := r.db.WithContext(ctx).Create(gu).Error; err != nil {
		// gorm ErrDuplicatedKey is driver-specific; you can check it too
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.logger.Warn().Err(err).Str("email", u.Email).Msg("duplicated key")
			return nil, errs.Wrap(err, errs.CodeConflict, "duplicated key")
		}

		r.logger.Error().Err(err).Str("email", u.Email).Msg("create user failed")
		return nil, errs.Wrap(err, errs.CodeInternal, "db create failed")
	}
	r.logger.Info().Uint64("user_id", gu.ID).Str("email", gu.Email).Msg("user created")
	return gu.ToEntity(), nil
}

func (r *userRepoImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var gu GormUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug().Str("email", email).Msg("user not found")
			return nil, errs.Wrap(err, errs.CodeNotFound, "user not found")
		}
		r.logger.Error().Err(err).Str("email", email).Msg("db query failed")
		return nil, errs.Wrap(err, errs.CodeInternal, "db query failed")
	}
	return gu.ToEntity(), nil
}

// TODO: update GetByID logging and error handling as Create and GetByEmail
func (r *userRepoImpl) GetByID(ctx context.Context, id uint64) (*entities.User, error) {
	var gu GormUser
	if err := r.db.WithContext(ctx).First(&gu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return gu.ToEntity(), nil
}

// TODO: update Update logging and error handling as Create and GetByEmail
func (r *userRepoImpl) Update(ctx context.Context, user *entities.User) error {
	gu := FromEntity(user)
	return r.db.WithContext(ctx).Save(gu).Error
}
