package ports

import (
	"context"
	"time"

	"github.com/itdyaingenieria/otp-service/internal/domain"

	"github.com/google/uuid"
)

type OTPRepository interface {
	Create(ctx context.Context, otp domain.OTPCode) error
	Get(ctx context.Context, tenantID string, id uuid.UUID) (domain.OTPCode, error)
	UpdateAttempts(ctx context.Context, tenantID string, id uuid.UUID, attempts int) error
	MarkUsed(ctx context.Context, tenantID string, id uuid.UUID, usedAt time.Time) error
}

type Notifier interface {
	Send(ctx context.Context, ch domain.Channel, to, message string) error
}

type CodeGenerator interface {
	Generate() (string, error)
}

type Clock interface{ Now() time.Time }
