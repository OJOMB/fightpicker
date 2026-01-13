package contextual

import (
	"context"
	"os/signal"
	"syscall"
)

func SetupSignals() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}
