package helpers

import "github.com/getsentry/sentry-go"

func SendTraceErrorToSentry(err error) error {
	if err != nil {
		sentry.CaptureException(err)
	}
	return err
}
