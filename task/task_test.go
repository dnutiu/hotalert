package task

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hotalert/alert"
	"testing"
	"time"
)

func Test_NewTask(t *testing.T) {
	var task = NewTask("web_scrape", Options{
		"option": "true",
	}, alert.NewDummyAlerter())
	assert.Equal(t, Task{
		ExecutionFuncName: "web_scrape",
		Options: Options{
			"option": "true",
		},
		Timeout:  10 * time.Second,
		Alerter:  alert.NewDummyAlerter(),
		Callback: nil,
	}, *task)
}

func Test_NewResult(t *testing.T) {
	var task = NewTask("web_scrape", Options{
		"option": "true",
	}, alert.NewDummyAlerter())
	testError := errors.New("test error")
	var result = NewResult(task)
	result.SetError(testError)
	assert.Equal(t, Result{
		InitialTask: &Task{
			ExecutionFuncName: "web_scrape",
			Options: Options{
				"option": "true",
			},
			Timeout:  10 * time.Second,
			Alerter:  alert.NewDummyAlerter(),
			Callback: nil,
		},
		error: testError,
	}, *result)
}
