package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/itdyaingenieria/otp-service/internal/adapter/clock"
	"github.com/itdyaingenieria/otp-service/internal/adapter/codegen"
	"github.com/itdyaingenieria/otp-service/internal/adapter/notifier"
	pgrepo "github.com/itdyaingenieria/otp-service/internal/adapter/repository/postgres"
	"github.com/itdyaingenieria/otp-service/internal/config"
	"github.com/itdyaingenieria/otp-service/internal/httpapi"
	"github.com/itdyaingenieria/otp-service/internal/usecase"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	repo := pgrepo.NewOTPRepository(pool)
	n := notifier.NewLogNotifier()
	cg := codegen.NewNumeric(6)
	cl := clock.System{}

	gen := usecase.NewGenerateOTP(repo, n, cg, cl, time.Duration(cfg.OTPTTLSec)*time.Second, cfg.OTPMaxAttempts)
	val := usecase.NewValidateOTP(repo, cl)

	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	h := httpapi.NewHandler(gen, val)
	h.Register(e)

	// health
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	addr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("listening on %s", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
