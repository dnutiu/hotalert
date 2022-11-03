package workload

import (
	"github.com/stretchr/testify/assert"
	"hotalert/task"
	"testing"
	"time"
)

func Test_FromYamlContent(t *testing.T) {
	var fileContents = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

	currentWorkload, err := FromYamlContent([]byte(fileContents))

	// General tests
	assert.NoError(t, err)
	assert.Len(t, currentWorkload.tasksList, 2)
	assert.Len(t, currentWorkload.alerterMap, 1)

	// Test alerter.
	alerter := currentWorkload.alerterMap["webhook_discord"]
	for _, taskEntry := range currentWorkload.tasksList {
		assert.Equal(t, alerter, taskEntry.Alerter)
	}

	// Test timeout
	assert.Equal(t, 10*time.Second, currentWorkload.tasksList[0].Timeout)
	assert.Equal(t, 15*time.Second, currentWorkload.tasksList[1].Timeout)

	// Test Options
	assert.Equal(t, task.Options{
		"url":      "https://jobs.eu",
		"keywords": []any{"Software Engineer, Backend"},
	}, currentWorkload.tasksList[0].Options)
	assert.Equal(t, task.Options{
		"url":         "https://jobs.ro",
		"keywords":    []any{"Software Engineer, Front-End", "Software Architect"},
		"extra_int":   80,
		"extra_float": 80.2,
		"extra_bool":  true,
	}, currentWorkload.tasksList[1].Options)
}

func Test_GetTasks(t *testing.T) {
	var fileContents = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

	currentWorkload, err := FromYamlContent([]byte(fileContents))
	assert.NoError(t, err)

	tasks := currentWorkload.GetTasks()
	assert.Len(t, tasks, 2)
	assert.IsType(t, []*task.Task{}, tasks)
}

func Test_GetTasksLen(t *testing.T) {
	var fileContents = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

	currentWorkload, err := FromYamlContent([]byte(fileContents))
	assert.NoError(t, err)
	assert.Equal(t, 2, currentWorkload.GetTasksLen())
}

var testAlertsKeyNotExistsContents = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alertx:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testAlertsKeyNotMapType = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts: "I'm just a string please don't hurt me."
`

var testAlertsKeyDuplicated = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testAlertsUnknownAlerter = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
  - options:
      url: https://jobs.ro
      keywords: ["Software Engineer, Front-End", "Software Architect"]
      extra_int: 80
      extra_float: 80.2
      extra_bool: True
    timeout: 15
    alerter: "webhook_discord"
alerts:
  imcoolalerter:
    cool: true
`

var testTasksDoesNotExists = `
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testTasksIsNotAnArray = `
tasks: "Many, very much tasks."
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testTasksTaskIsNotAMap = `
tasks:
  - task: "cool task"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testTasksTaskHasInvalidAlerter = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "imaacoolalerter"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

var testTasksListEmpty = `
tasks: []
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

func Test_FromYamlContent_Errors(t *testing.T) {
	var tests = []struct {
		TestName         string
		TestFileContents string
	}{
		{
			"AlertsKeyDoesNotExists",
			testAlertsKeyNotExistsContents,
		},
		{
			"testAlertsKeyNotMapType",
			testAlertsKeyNotMapType,
		},
		{
			"testAlertsKeyDuplicated",
			testAlertsKeyDuplicated,
		},
		{
			"testAlertsUnknownAlerter",
			testAlertsUnknownAlerter,
		},
		{
			"testTasksDoesNotExists",
			testTasksDoesNotExists,
		},
		{
			"testTasksIsNotAnArray",
			testTasksIsNotAnArray,
		},
		{
			"testTasksTaskIsNotAMap",
			testTasksTaskIsNotAMap,
		},
		{
			"testTasksListEmpty",
			testTasksListEmpty,
		},
		{
			"testTasksTaskHasInvalidAlerter",
			testTasksTaskHasInvalidAlerter,
		},
	}

	for _, tv := range tests {
		t.Run(tv.TestName, func(t *testing.T) {
			currentWorkload, err := FromYamlContent([]byte(tv.TestFileContents))
			assert.Nil(t, currentWorkload)
			assert.Error(t, err)
		})
	}
}

var testTasksTaskHasInvalidAlerter2 = `
tasks:
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "imaacoolalerter"
  - options:
      url: https://jobs.eu
      keywords: ["Software Engineer, Backend"]
    timeout: 10
    alerter: "webhook_discord"
alerts:
  webhook_discord:
    webhook: https://webhook.url.com
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword $keywords was found on page!"
`

func Test_FromYamlContent_InvalidAlerterForTask(t *testing.T) {
	currentWorkload, err := FromYamlContent([]byte(testTasksTaskHasInvalidAlerter2))
	assert.NoError(t, err)
	assert.NotNil(t, currentWorkload)
	assert.Len(t, currentWorkload.tasksList, 1)
}
