package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	cfgpkg "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/config"
	dbpkg "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/db"
	"github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/jwt"
	logpkg "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/log"

	dbadapter "github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/adapters/db"
	"github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/adapters/http"
	httpadapter "github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/adapters/http"
	usecases "github.com/BitKa-Exchange/bitka-exchange/auth-service/internal/usecases"
	jwtpkg "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/jwt"
	passhash "github.com/BitKa-Exchange/bitka-exchange/auth-service/pkg/passhash"
)

func main() {
	cfg, err := cfgpkg.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger := logpkg.NewLogger(cfg)
	logger.Info().Msg("config loaded")

	gormDB, err := dbpkg.NewGorm(cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect db")
		os.Exit(1)
	}

	userRepo := dbadapter.NewUserRepo(gormDB, logger) // implements usecases.UserRepository
	// TODO: Complete RefreshTokenRepository implementation
	jwtSvc := jwtpkg.NewHMACService(jwt.HMACConfig{
		AccessSecret: []byte(cfg.JwtAccessSecret),
		AccessTTL:    cfg.AccessTTL,
		RefreshTTL:   cfg.RefreshTTL,
		Issuer:       "fill here", // e.g. your service name
		Audience:     "fill here", // e.g. your service audience
	})
	hasher := passhash.NewBcrypt(12)

	authUC := usecases.NewAuthUsecase(userRepo, jwtSvc, hasher, logger) // pass logger if needed

	app := fiber.New(fiber.Config{
		ErrorHandler: http.ErrorHandler(logger),
	})

	h := httpadapter.NewHandler(authUC)
	h.RegisterRoutes(app)

	port := cfg.Port
	logger.Info().Str("port", port).Msg("starting http server")
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal().Err(err).Msg("server stopped")
	}
}
