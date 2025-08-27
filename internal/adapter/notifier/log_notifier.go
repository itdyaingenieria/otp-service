package notifier

import (
	"context"
	"log"

	"github.com/itdyaingenieria/otp-service/internal/domain"
)

type LogNotifier struct{}

func NewLogNotifier() LogNotifier { return LogNotifier{} }

func (LogNotifier) Send(ctx context.Context, ch domain.Channel, to, msg string) error {
	log.Printf("[NOTIFIER] channel=%s to=%s msg=%s", ch, to, msg)
	return nil
}
