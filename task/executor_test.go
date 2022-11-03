package task

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hotalert/alert"
	"testing"
)

func Test_DefaultExecutor(t *testing.T) {
	// Setup
	var taskCounter = 0
	defaultExecutor := NewDefaultExecutor(func(task *Task) error {
		// First task is successful, others return error.
		if taskCounter > 0 {
			return errors.New("test")
		}
		taskCounter += 1
		return nil
	})
	taskResultsChan := defaultExecutor.Start()

	// Test
	var task1 = &Task{
		Timeout:  0,
		Alerter:  alert.NewDummyAlerter(),
		Callback: nil,
	}
	var task2 = &Task{
		Timeout:  0,
		Alerter:  alert.NewDummyAlerter(),
		Callback: nil,
	}

	defaultExecutor.AddTask(task1)
	defaultExecutor.AddTask(task2)

	// Assert results
	assert.Equal(t, &Result{
		InitialTask: task1,
		error:       nil,
	}, <-taskResultsChan)
	assert.Equal(t, &Result{
		InitialTask: task2,
		error:       errors.New("test"),
	}, <-taskResultsChan)

	// Clean-up
	defaultExecutor.Shutdown()
}
