package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/itdyaingenieria/otp-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type otpRepository struct{ pool *pgxpool.Pool }

func NewOTPRepository(pool *pgxpool.Pool) *otpRepository { return &otpRepository{pool: pool} }

func (r *otpRepository) Create(ctx context.Context, o domain.OTPCode) error {
	_, err := r.pool.Exec(ctx, `
INSERT INTO otps (id, tenant_id, channel, destination, code, attempts, max_attempts, expires_at, used_at, created_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
`, o.ID, o.TenantID, string(o.Channel), o.Destination, o.Code, o.Attempts, o.MaxAttempts, o.ExpiresAt, o.UsedAt, o.CreatedAt)
	return err
}

func (r *otpRepository) Get(ctx context.Context, tenantID string, id uuid.UUID) (domain.OTPCode, error) {
	row := r.pool.QueryRow(ctx, `
SELECT id, tenant_id, channel, destination, code, attempts, max_attempts, expires_at, used_at, created_at
FROM otps WHERE id=$1 AND tenant_id=$2
`, id, tenantID)
	var o domain.OTPCode
	var ch string
	if err := row.Scan(&o.ID, &o.TenantID, &ch, &o.Destination, &o.Code, &o.Attempts, &o.MaxAttempts, &o.ExpiresAt, &o.UsedAt, &o.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OTPCode{}, domain.ErrNotFound
		}
		return domain.OTPCode{}, err
	}
	o.Channel = domain.Channel(ch)
	return o, nil
}

func (r *otpRepository) UpdateAttempts(ctx context.Context, tenantID string, id uuid.UUID, attempts int) error {
	ct, err := r.pool.Exec(ctx, `
UPDATE otps SET attempts=$1 WHERE id=$2 AND tenant_id=$3
`, attempts, id, tenantID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *otpRepository) MarkUsed(ctx context.Context, tenantID string, id uuid.UUID, usedAt time.Time) error {
	ct, err := r.pool.Exec(ctx, `
UPDATE otps SET used_at=$1 WHERE id=$2 AND tenant_id=$3
`, usedAt, id, tenantID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
