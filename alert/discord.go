package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hotalert/logging"
	"net/http"
	"strings"
)

// DiscordWebhookAlerter is a struct that implements alerting on Discord via webhooks.
type DiscordWebhookAlerter struct {
	// webhook is the Discord webhook URL.
	webhook string
	// messageTemplate is the message that is going to be posted when the alert conditions match.
	messageTemplate string
	// HttpClient is the http client used when executing requests.
	HttpClient *http.Client
}

// DiscordWebhookAlerterOptions are the options for the DiscordWebhookAlerter
type DiscordWebhookAlerterOptions struct {
	// Webhook is the discord webhook.
	Webhook string `mapstructure:"webhook"`
	// MessageTemplate is the message template that is going to be posted.
	MessageTemplate string `mapstructure:"message"`
}

// Validate validates the DiscordWebhookAlerterOptions, returns an error on invalid options.
func (o *DiscordWebhookAlerterOptions) Validate() error {
	if o.Webhook == "" || o.MessageTemplate == "" {
		return errors.New("invalid configuration for webhook_discord")
	}
	if !strings.Contains(o.Webhook, "http://") && !strings.Contains(o.Webhook, "https://") {
		return errors.New(fmt.Sprintf("invalid webhook schema for %s", o.Webhook))
	}
	return nil
}

// NewDiscordWebhookAlerter returns a new DiscordWebhookAlerter instance.
func NewDiscordWebhookAlerter(options DiscordWebhookAlerterOptions) (*DiscordWebhookAlerter, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	return &DiscordWebhookAlerter{
		webhook:         options.Webhook,
		messageTemplate: options.MessageTemplate,
		HttpClient:      http.DefaultClient,
	}, nil
}

// PostAlert posts the alert on Discord via webhooks.
func (d *DiscordWebhookAlerter) PostAlert(ctx context.Context, matchedKeywords []string) {
	alertMessage := strings.Replace(d.messageTemplate, "$keywords", strings.Join(matchedKeywords, ","), -1)
	var postBody = map[string]interface{}{
		"content":     alertMessage,
		"embeds":      nil,
		"attachments": nil,
	}

	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		logging.SugaredLogger.Errorf("Failed to marshall postBody: %v", err)
		return
	}

	postRequest, err := http.NewRequest("POST", d.webhook, bytes.NewBuffer(postBodyBytes))
	if err != nil {
		logging.SugaredLogger.Errorf("Failed to create alert request.")
		return
	}
	postRequest.Header["Content-Type"] = []string{"application/json"}
	_, err = d.HttpClient.Do(postRequest.WithContext(ctx))
	if err != nil {
		logging.SugaredLogger.Errorf("Failed to post alert to discord!")
		return
	}
	logging.SugaredLogger.Infof("Alert posted:\nBEGIN\n%s\nEND", alertMessage)
}
