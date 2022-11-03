package alert

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_DiscordWebhookAlerterOptions_Validate(t *testing.T) {
	var tests = []struct {
		Webhook         string
		MessageTemplate string
		IsValid         bool
	}{
		{
			"",
			"",
			false,
		},
		{
			"",
			"asdasd",
			false,
		},
		{
			"asdasd",
			"asdasd",
			false,
		},
		{
			"http://example.com",
			"The template",
			true,
		},
		{
			"https://example.com",
			"The template",
			true,
		},
	}

	for ti, tv := range tests {
		t.Run(fmt.Sprintf("test_%d", ti), func(t *testing.T) {
			opts := DiscordWebhookAlerterOptions{Webhook: tv.Webhook, MessageTemplate: tv.MessageTemplate}
			err := opts.Validate()
			if tv.IsValid {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_DiscordWebhookAlerter_PostAlert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"attachments\":null,\"content\":\"test matched,second\",\"embeds\":null}", string(requestBody))
		assert.Equal(t, "application/json", r.Header["Content-Type"][0])
	}))
	defer ts.Close()

	client := ts.Client()
	alerter, err := NewDiscordWebhookAlerter(DiscordWebhookAlerterOptions{
		Webhook:         ts.URL,
		MessageTemplate: "test $keywords",
	})
	assert.NoError(t, err)
	alerter.HttpClient = client

	alerter.PostAlert(context.Background(), []string{"matched", "second"})
}
