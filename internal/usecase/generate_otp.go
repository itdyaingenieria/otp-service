package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/itdyaingenieria/otp-service/internal/domain"
	"github.com/itdyaingenieria/otp-service/internal/ports"

	"github.com/google/uuid"
)

type GenerateOTP struct {
	repo       ports.OTPRepository
	n          ports.Notifier
	cg         ports.CodeGenerator
	clock      ports.Clock
	ttl        time.Duration
	maxAttempt int
}

func NewGenerateOTP(repo ports.OTPRepository, n ports.Notifier, cg ports.CodeGenerator, clock ports.Clock, ttl time.Duration, maxAttempt int) GenerateOTP {
	return GenerateOTP{repo: repo, n: n, cg: cg, clock: clock, ttl: ttl, maxAttempt: maxAttempt}
}

type GenerateInput struct {
	TenantID    string
	Channel     domain.Channel
	Destination string
}

type GenerateOutput struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (uc GenerateOTP) Execute(ctx context.Context, in GenerateInput) (GenerateOutput, error) {
	code, err := uc.cg.Generate()
	if err != nil {
		return GenerateOutput{}, err
	}
	now := uc.clock.Now()
	otp := domain.OTPCode{
		ID:          uuid.New(),
		TenantID:    in.TenantID,
		Channel:     in.Channel,
		Destination: in.Destination,
		Code:        code,
		Attempts:    0,
		MaxAttempts: uc.maxAttempt,
		ExpiresAt:   now.Add(uc.ttl),
		CreatedAt:   now,
	}
	if err := uc.repo.Create(ctx, otp); err != nil {
		return GenerateOutput{}, err
	}
	msg := fmt.Sprintf("Your verification code is %s. It expires at %s.", code, otp.ExpiresAt.Format(time.RFC3339))
	_ = uc.n.Send(ctx, in.Channel, in.Destination, msg) // ignore notifier error for demo; in prod handle & retry
	return GenerateOutput{ID: otp.ID, ExpiresAt: otp.ExpiresAt}, nil
}
