package alert

import (
	"errors"
	"fmt"
)

// NewAlerter builds Alerted function given the alerter name and options.
func NewAlerter(name string, options map[string]interface{}) (Alerter, error) {
	if name == "webhook_discord" {
		return NewDiscordWebhookAlerter(DiscordWebhookAlerterOptions{
			Webhook:         options["webhook"].(string),
			MessageTemplate: options["message"].(string),
		})
	}
	return nil, errors.New(fmt.Sprintf("invalid alerter name %s", name))
}
