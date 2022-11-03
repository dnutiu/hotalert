package alert

import (
	"context"
	"testing"
)

func TestDummyAlerter_PostAlert(t *testing.T) {
	NewDummyAlerter().PostAlert(context.TODO(), []string{"demo"})
}
