package users

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	EmailVerificationAttempts metric.Int64Counter
	EmailVerificationFailures metric.Int64Counter
}

func New() (*Metrics, error) {
	meter := otel.Meter("fightpicker.users")

	emailVerificationAttempts, err := meter.Int64Counter(
		"user.email_verification.attempts",
		metric.WithDescription("Number of email verification attempts"),
	)
	if err != nil {
		return nil, err
	}

	emailVerificationFailures, err := meter.Int64Counter(
		"user.email_verification.failures",
		metric.WithDescription("Number of failed email verification attempts"),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{
		EmailVerificationAttempts: emailVerificationAttempts,
		EmailVerificationFailures: emailVerificationFailures,
	}, nil
}
