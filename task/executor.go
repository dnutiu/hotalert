package task

import (
	"errors"
	"fmt"
	"sync"
)

// Executor is an interface for implementing task executors.
type Executor interface {
	// AddTask adds a task to the task executor.
	AddTask(task *Task)
	// Start starts the scrapper and returns a Result receive-only channel.
	Start() <-chan *Result
	// Shutdown shuts down the scrapper. It will block until the Executor was shut down.
	Shutdown()
}

// ExecutionFn is a type definition for a function that executes the task and returns an error.
type ExecutionFn func(task *Task) error

// DefaultExecutor is a TaskExecutor with the default implementation.
// The tasks are executed directly on the machine.
type DefaultExecutor struct {
	// TaskExecutionFn is the function that executes the task.
	TaskExecutionFn ExecutionFn
	// workerGroup is a waiting group for worker goroutines.
	workerGroup *sync.WaitGroup
	// numberOfWorkerGoroutines is the number of working goroutines.
	numberOfWorkerGoroutines int
	// taskResultChan is a receive only channel for task results.
	taskResultChan chan *Result
	// taskChan is a channel for tasks
	taskChan chan *Task
	// quinChan is a channel for sending the quit command to worker goroutines.
	quinChan chan int
}

// TODO: Add support for per task execution functions

// NewDefaultExecutor returns a new instance of DefaultExecutor.
func NewDefaultExecutor(fn ExecutionFn) *DefaultExecutor {
	ws := &DefaultExecutor{
		TaskExecutionFn:          fn,
		workerGroup:              &sync.WaitGroup{},
		taskResultChan:           make(chan *Result, 50),
		taskChan:                 make(chan *Task, 50),
		numberOfWorkerGoroutines: 5,
	}
	ws.quinChan = make(chan int, ws.numberOfWorkerGoroutines)
	return ws
}

// AddTask adds a task to the DefaultExecutor queue.
func (ws *DefaultExecutor) AddTask(task *Task) {
	ws.taskChan <- task
}

// executeTask executes the given task using TaskExecutionFn
func (ws *DefaultExecutor) executeTask(task *Task) error {
	var taskErr error = nil
	// Execute task and set panics as errors in taskResult.
	func() {
		defer func() {
			if r := recover(); r != nil {
				taskErr = errors.New(fmt.Sprintf("panic: %s", r))
			}
		}()
		err := ws.TaskExecutionFn(task)
		if err != nil {
			taskErr = err
		}
	}()
	return taskErr
}

// workerGoroutine waits for tasks and executes them.
// After the task is executed it forwards the result, including errors and panics to the task Result channel.
func (ws *DefaultExecutor) workerGoroutine() {
	defer ws.workerGroup.Done()
	for {
		select {
		case task := <-ws.taskChan:
			var taskResult = NewResult(task)
			taskResult.error = ws.executeTask(task)

			// Forward TaskResult to channel.
			ws.taskResultChan <- taskResult
		case <-ws.quinChan:
			// Quit
			return
		}
	}
}

// Start starts the DefaultExecutor.
// Start returns a receive only channel with task Result.
func (ws *DefaultExecutor) Start() <-chan *Result {
	// Start worker goroutines.
	for i := 0; i < ws.numberOfWorkerGoroutines; i++ {
		ws.workerGroup.Add(1)
		go ws.workerGoroutine()
	}
	return ws.taskResultChan
}

// Shutdown shuts down the DefaultExecutor.
// Shutdown blocks till the DefaultExecutor has shutdown.
func (ws *DefaultExecutor) Shutdown() {
	// Shutdown all worker goroutines
	for i := 0; i < ws.numberOfWorkerGoroutines; i++ {
		ws.quinChan <- 1
	}
	ws.workerGroup.Wait()
	close(ws.taskChan)
	close(ws.taskResultChan)
	close(ws.quinChan)
}
