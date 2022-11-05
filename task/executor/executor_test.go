package executor

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/stretchr/testify/assert"
	"hotalert/alert"
	"hotalert/task"
	"testing"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Test_DefaultExecutor(t *testing.T) {
	// Setup
	var taskCounter = 0
	var taskTestFunc = func(task *task.Task) error {
		// First task is successful, others return error.
		if taskCounter > 0 {
			return errors.New("test")
		}
		taskCounter += 1
		return nil
	}

	err := RegisterNewExecutionFunction("task_test", taskTestFunc)
	assert.NoError(t, err)

	defaultExecutor := NewDefaultExecutor()
	taskResultsChan := defaultExecutor.Start()

	// Test
	var task1 = &task.Task{
		ExecutionFuncName: "task_test",
		Timeout:           0,
		Alerter:           alert.NewDummyAlerter(),
		Callback:          nil,
	}
	var task2 = &task.Task{
		ExecutionFuncName: "task_test",
		Timeout:           0,
		Alerter:           alert.NewDummyAlerter(),
		Callback:          nil,
	}

	defaultExecutor.AddTask(task1)
	defaultExecutor.AddTask(task2)

	// Assert results
	result1 := <-taskResultsChan
	assert.Equal(t, task1, result1.InitialTask)
	assert.Equal(t, nil, result1.Error())

	result2 := <-taskResultsChan
	assert.Equal(t, task2, result2.InitialTask)
	assert.Equal(t, errors.New("test"), result2.Error())

	// Clean-up
	defaultExecutor.Shutdown()
}

func Test_DefaultExecutor_InvalidTaskExecutionFuncName(t *testing.T) {
	defaultExecutor := NewDefaultExecutor()
	taskResultsChan := defaultExecutor.Start()

	// Test
	var task1 = &task.Task{
		ExecutionFuncName: "vand_dacia_2006",
		Timeout:           0,
		Alerter:           alert.NewDummyAlerter(),
		Callback:          nil,
	}
	defaultExecutor.AddTask(task1)

	// Assert results
	result1 := <-taskResultsChan
	assert.Equal(t, task1, result1.InitialTask)
	assert.Equal(t, errors.New("invalid task execution function name: 'vand_dacia_2006'"), result1.Error())

	// Clean-up
	defaultExecutor.Shutdown()
}

func Test_RegisterNewExecutionFunction(t *testing.T) {
	var taskTestFunc = func(t *task.Task) error { return nil }
	randomName, _ := randomHex(5)
	err := RegisterNewExecutionFunction(randomName, taskTestFunc)
	assert.NoError(t, err)
	err = RegisterNewExecutionFunction(randomName, taskTestFunc)
	assert.Error(t, err)
}
