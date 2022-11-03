package alert

import (
	"context"
	"hotalert/logging"
)

// DummyAlerter is an Alerter that does nothing. It is used when no actual alerter is available.
type DummyAlerter struct {
}

// NewDummyAlerter returns a new instance of DummyAlerter.
func NewDummyAlerter() *DummyAlerter {
	return &DummyAlerter{}
}

func (d DummyAlerter) PostAlert(ctx context.Context, matchedKeywords []string) {
	logging.SugaredLogger.Infof("DummyAlert: %v - %v", ctx, matchedKeywords)
}
