package config

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func InitSentry(cfg *Config) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Sentry.Dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		Debug:            true,
	})
	if err != nil {
		return err
	}

	return nil
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
	log.Println("Sentry flushed.")
}
