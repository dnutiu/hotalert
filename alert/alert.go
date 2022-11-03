package alert

import "context"

// Alerter is an interface for implementing alerts on various channels
type Alerter interface {
	// PostAlert posts the alert with the given message.
	PostAlert(ctx context.Context, matchedKeywords []string)
}
