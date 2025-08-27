package usecase

import (
	"context"

	"github.com/itdyaingenieria/otp-service/internal/domain"
	"github.com/itdyaingenieria/otp-service/internal/ports"

	"github.com/google/uuid"
)

type ValidateOTP struct {
	repo  ports.OTPRepository
	clock ports.Clock
}

func NewValidateOTP(repo ports.OTPRepository, clock ports.Clock) ValidateOTP {
	return ValidateOTP{repo: repo, clock: clock}
}

type ValidateInput struct {
	TenantID string
	ID       uuid.UUID
	Code     string
}

type ValidateOutput struct {
	Valid bool `json:"valid"`
}

func (uc ValidateOTP) Execute(ctx context.Context, in ValidateInput) (ValidateOutput, error) {
	otp, err := uc.repo.Get(ctx, in.TenantID, in.ID)
	if err != nil {
		return ValidateOutput{}, err
	}
	if otp.IsUsed() {
		return ValidateOutput{}, domain.ErrAlreadyUsed
	}
	if otp.IsExpired(uc.clock.Now()) {
		return ValidateOutput{}, domain.ErrExpired
	}
	if !otp.CanAttempt() {
		return ValidateOutput{}, domain.ErrMaxAttempts
	}

	if in.Code != otp.Code {
		// wrong attempt
		if err := uc.repo.UpdateAttempts(ctx, in.TenantID, in.ID, otp.Attempts+1); err != nil {
			return ValidateOutput{}, err
		}
		return ValidateOutput{Valid: false}, domain.ErrInvalidCode
	}
	// success
	now := uc.clock.Now()
	if err := uc.repo.MarkUsed(ctx, in.TenantID, in.ID, now); err != nil {
		return ValidateOutput{}, err
	}
	return ValidateOutput{Valid: true}, nil
}
