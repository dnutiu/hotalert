package task

import (
	"fmt"
	"hotalert/alert"
	"time"
)

// Callback represents a callback function that is called after task completion
type Callback func(result *Result)

// Options represents the options available for a task.
type Options map[string]any

// Task represents the context of a task.
type Task struct {
	// ExecutionFuncName is the function name associated with this task.
	ExecutionFuncName string
	// Options are the option given to the task.
	Options Options `mapstructure:"options"`
	// Timeout is the timeout for the task.
	Timeout time.Duration `mapstructure:"timeout"`
	// Alerter is the alerter that will be called when task is completed.
	Alerter alert.Alerter `mapstructure:"alerter"`
	// Callback is an optional function that will be called when task is completed. (Not implemented)
	Callback *Callback
}

// NewTask returns a new task instance.
func NewTask(executionFuncName string, options Options, alerter alert.Alerter) *Task {
	if alerter == nil {
		panic(fmt.Sprintf("Alerter cannot be nil"))
	}
	return &Task{
		ExecutionFuncName: executionFuncName,
		Options:           options,
		Timeout:           10 * time.Second,
		Alerter:           alerter,
		Callback:          nil,
	}
}

// Result represents the result of a task.
type Result struct {
	// InitialTask is the original Task for which the Result is given.
	InitialTask *Task
	// error is the error of the task.
	error error
}

// NewResult represents the result of a task.
func NewResult(task *Task) *Result {
	return &Result{
		InitialTask: task,
		error:       nil,
	}
}

// SetError sets the error on the result object.
func (r *Result) SetError(err error) {
	r.error = err
}

// Error returns the error encountered during the execution of the task.
// Error returns null if the task had no errors and was completed.
func (r *Result) Error() error {
	return r.error
}
