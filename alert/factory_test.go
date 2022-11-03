package alert

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewAlerter(t *testing.T) {
	var tests = []struct {
		TestName       string
		AlerterName    string
		AlerterOptions map[string]interface{}
		ExpectedType   interface{}
		ShouldError    bool
	}{
		{
			TestName:    "Webhook Discord",
			AlerterName: "webhook_discord",
			AlerterOptions: map[string]interface{}{
				"webhook": "https://webhook.test",
				"message": "The Message is fine.",
			},
			ExpectedType: &DiscordWebhookAlerter{},
			ShouldError:  false,
		},
		{
			TestName:    "Webhook Discord Error",
			AlerterName: "webhook_discord",
			AlerterOptions: map[string]interface{}{
				"webhook": "",
				"message": "The Message is fine.",
			},
			ExpectedType: &DiscordWebhookAlerter{},
			ShouldError:  true,
		},
	}

	for _, tv := range tests {
		t.Run(tv.TestName, func(t *testing.T) {
			alerter, err := NewAlerter(tv.AlerterName, tv.AlerterOptions)
			if !tv.ShouldError {
				assert.Nil(t, err)
				assert.IsType(t, tv.ExpectedType, alerter)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, alerter)
			}
		})
	}
}
