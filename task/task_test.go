package task

import (
	"github.com/stretchr/testify/assert"
	"hotalert/alert"
	"testing"
	"time"
)

func Test_NewTask(t *testing.T) {
	var task = NewTask(Options{
		"option": "true",
	}, alert.NewDummyAlerter())
	assert.Equal(t, Task{
		Options: Options{
			"option": "true",
		},
		Timeout:  10 * time.Second,
		Alerter:  alert.NewDummyAlerter(),
		Callback: nil,
	}, *task)
}

func Test_NewResult(t *testing.T) {
	var task = NewTask(Options{
		"option": "true",
	}, alert.NewDummyAlerter())
	var result = NewResult(task)
	assert.Equal(t, Result{
		InitialTask: &Task{
			Options: Options{
				"option": "true",
			},
			Timeout:  10 * time.Second,
			Alerter:  alert.NewDummyAlerter(),
			Callback: nil,
		},
		error: nil,
	}, *result)
}
