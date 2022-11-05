package executor

import (
	"errors"
	"fmt"
	"hotalert/logging"
	"hotalert/task"
	"hotalert/task/functions"
	"sync"
)

// Executor is an interface for implementing task executors.
type Executor interface {
	// AddTask adds a task to the task executor.
	AddTask(task *task.Task)
	// Start starts the scrapper and returns a Result receive-only channel.
	Start() <-chan *task.Result
	// Shutdown shuts down the scrapper. It will block until the Executor was shut down.
	Shutdown()
}

// ExecutionFunc is a type definition for a function that executes the task and returns an error.
type ExecutionFunc func(task *task.Task) error

// DefaultExecutor is a TaskExecutor with the default implementation.
// The tasks are executed directly on the machine.
type DefaultExecutor struct {
	// workerGroup is a waiting group for worker goroutines.
	workerGroup *sync.WaitGroup
	// numberOfWorkerGoroutines is the number of working goroutines.
	numberOfWorkerGoroutines int
	// taskResultChan is a receive only channel for task results.
	taskResultChan chan *task.Result
	// taskChan is a channel for tasks
	taskChan chan *task.Task
	// quinChan is a channel for sending the quit command to worker goroutines.
	quinChan chan int
}

// executionFuncMap is a map that holds all the possible values for ExecutionFunc.
// Right now it is hard-coded but in the future it may be extended dynamically.
var executionFuncMap = map[string]ExecutionFunc{
	"web_scrape": functions.WebScrapeTask,
}

// RegisterNewExecutionFunction registers a new execution function.
func RegisterNewExecutionFunction(name string, function ExecutionFunc) error {
	for n := range executionFuncMap {
		if n == name {
			return errors.New("function already exists")
		}
	}
	executionFuncMap[name] = function
	return nil
}

// NewDefaultExecutor returns a new instance of DefaultExecutor.
func NewDefaultExecutor() *DefaultExecutor {
	ws := &DefaultExecutor{
		workerGroup:              &sync.WaitGroup{},
		taskResultChan:           make(chan *task.Result, 50),
		taskChan:                 make(chan *task.Task, 50),
		numberOfWorkerGoroutines: 5,
	}
	ws.quinChan = make(chan int, ws.numberOfWorkerGoroutines)
	return ws
}

// AddTask adds a task to the DefaultExecutor queue.
func (ws *DefaultExecutor) AddTask(task *task.Task) {
	ws.taskChan <- task
}

// executeTask executes the given task using DefaultTaskExecutionFuncName
func (ws *DefaultExecutor) executeTask(task *task.Task) error {
	var taskErr error = nil
	// Execute task and set panics as errors in taskResult.
	func() {
		defer func() {
			if r := recover(); r != nil {
				taskErr = errors.New(fmt.Sprintf("panic: %s", r))
			}
		}()

		taskExecutionFunc, ok := executionFuncMap[task.ExecutionFuncName]
		if !ok {
			message := fmt.Sprintf("invalid task execution function name: '%s'", task.ExecutionFuncName)
			logging.SugaredLogger.Error(message)
			taskErr = errors.New(message)
			return
		}

		err := taskExecutionFunc(task)
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
		case currentTask := <-ws.taskChan:
			var taskResult = task.NewResult(currentTask)
			err := ws.executeTask(currentTask)
			taskResult.SetError(err)

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
func (ws *DefaultExecutor) Start() <-chan *task.Result {
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
