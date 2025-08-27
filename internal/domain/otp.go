package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
)

var (
	ErrNotFound    = errors.New("otp not found")
	ErrExpired     = errors.New("otp expired")
	ErrMaxAttempts = errors.New("max attempts reached")
	ErrAlreadyUsed = errors.New("otp already used")
	ErrInvalidCode = errors.New("invalid code")
)

type OTPCode struct {
	ID          uuid.UUID
	TenantID    string
	Channel     Channel
	Destination string
	Code        string
	Attempts    int
	MaxAttempts int
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   time.Time
}

func (o OTPCode) IsExpired(now time.Time) bool { return now.After(o.ExpiresAt) }
func (o OTPCode) IsUsed() bool                 { return o.UsedAt != nil }
func (o OTPCode) CanAttempt() bool             { return o.Attempts < o.MaxAttempts }
